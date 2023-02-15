package pubsub

import (
	"context"

	"github.com/mediocregopher/radix/v3"
)

// PubSub is a client for recieving messages using redis pubsub
type PubSub struct {
	conn   radix.PubSubConn
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new PubSub client and establishes the connection to redis
func New(address string, password string, tls bool) (*PubSub, error) {
	var dialOpts []radix.DialOpt
	dialOpts = append(dialOpts, radix.DialAuthPass(password))
	if tls {
		dialOpts = append(dialOpts, radix.DialUseTLS(nil))
	}
	connFunc := radix.PersistentPubSubConnFunc(func(string, address string) (radix.Conn, error) {
		return radix.Dial("tcp", address, dialOpts...)
	})

	conn, err := radix.PersistentPubSubWithOpts("tcp", address, connFunc, radix.PersistentPubSubAbortAfter(3))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &PubSub{
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Subscribe subscribes to a redis pubsub channel, and returns a channel for receiving messages
func (p *PubSub) Subscribe(channel string) (<-chan []byte, error) {
	in := make(chan radix.PubSubMessage)

	err := p.conn.Subscribe(in, channel)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)
	go p.worker(channel, in, out)

	return out, nil
}

func (p *PubSub) worker(channel string, in chan radix.PubSubMessage, out chan<- []byte) {
	defer p.cleanup(channel, in, out)

	for {
		select {
		case msg, open := <-in:
			if !open {
				return
			}

			out <- msg.Message
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *PubSub) cleanup(channel string, in chan radix.PubSubMessage, out chan<- []byte) {
	p.conn.Unsubscribe(in, channel)
	close(in)
	close(out)
}

// Shutdown shuts everything down and closes the redis connection
func (p *PubSub) Shutdown() {
	p.cancel()
	p.conn.Close()
}
