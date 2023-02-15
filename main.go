package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/infosum/statsd"
	"github.com/mullvad/message-queue/pubsub"
	"github.com/mullvad/message-queue/queue"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jamiealquiza/envy"
	"github.com/mullvad/message-queue/api"
)

// Build information populated as build-time.
var (
	Version   string
	Branch    string
	Revision  string
	GoVersion = runtime.Version()
)

// init creates a metrics about current version information.
func init() {
	promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "messagequeue",
			Name:      "build_info",
			Help:      "A metric with a constant '1' value labeled by version, branch, revision and goversion from which message-queue was built.",
			ConstLabels: prometheus.Labels{
				"version":   Version,
				"branch":    Branch,
				"revision":  Revision,
				"goversion": GoVersion,
			},
		},
	).Set(1)
}

var (
	metrics *statsd.Client
	p       *pubsub.PubSub
	q       *queue.Queue
)

func main() {
	// Set up commandline flags
	listen := flag.String("listen", ":8080", "listen address")
	bufferSize := flag.Int("buffer-size", 100, "client buffer size")
	redisAddress := flag.String("redis-address", "", "address of the redis server")
	redisPassword := flag.String("redis-password", "", "password for the redis server")
	channels := flag.String("channels", "", "comma-delimited list of channels to listen and broadcast to")
	statsdAddress := flag.String("statsd-address", "127.0.0.1:8125", "statsd address to send metrics to")

	// Parse environment variables
	envy.Parse("MQ")

	// Parse commandline flags
	flag.Parse()

	if *channels == "" {
		log.Fatalf("no channels configured")
	}

	channelList := strings.Split(*channels, ",")

	log.Printf("starting message-queue")

	// Initialize metrics
	var err error
	metrics, err = statsd.New(statsd.TagsFormat(statsd.Datadog), statsd.Prefix("mq"), statsd.Address(*statsdAddress))
	if err != nil {
		log.Fatal("Error initializing metrics: ", err)
	}
	defer metrics.Close()

	// Set up context for shutting down
	shutdownCtx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	// Set up the pubsub listener
	p, err = pubsub.New(*redisAddress, *redisPassword)
	if err != nil {
		log.Fatal("error initializing pubsub: ", err)
	}

	// Set up the queue
	q = queue.New(shutdownCtx, *bufferSize)

	// Set up the message passing from redis pubsub to the queue
	err = setupChannels(shutdownCtx, channelList)
	if err != nil {
		log.Fatal("error initializing queue: ", err)
	}

	// Create a ticker for metrics
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				collectMetrics()
			case <-shutdownCtx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	// Expose prometheus metrics
	log.Printf("Exposing metrics on port 9999")
	startMetricsServer(":9999")

	// Start and listen on http
	api := api.New(q)

	server := &http.Server{
		Addr:              *listen,
		Handler:           api.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("shutting down the http server", err)
			shutdown()
		}
	}()

	log.Printf("http server listening on %s", *listen)

	// Wait for shutdown or error
	err = waitForInterrupt(shutdownCtx)
	log.Println("shutting down", err)

	// Shut down http server
	serverCtx, serverCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer serverCancel()
	if err := server.Shutdown(serverCtx); err != nil {
		log.Println("error shutting down", err)
	}
}

func setupChannels(ctx context.Context, channels []string) error {
	for _, channel := range channels {
		in, err := p.Subscribe(channel)
		if err != nil {
			return err
		}

		out, err := q.CreateChannel(channel)
		if err != nil {
			return err
		}

		go channelWorker(ctx, in, out)
	}

	return nil
}

func channelWorker(ctx context.Context, in <-chan []byte, out chan<- []byte) {
	defer func() {
		close(out)
	}()

	for {
		select {
		case msg, open := <-in:
			if !open {
				return
			}

			out <- msg
		case <-ctx.Done():
			return
		}
	}
}

func collectMetrics() {
	metrics.Gauge("subscribers", q.SubscriberCount())
}

func startMetricsServer(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() { log.Fatal(server.ListenAndServe()) }()
}

func waitForInterrupt(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		return fmt.Errorf("received signal %s", sig)
	case <-ctx.Done():
		return errors.New("canceled")
	}
}
