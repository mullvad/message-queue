package pubsub_test

import (
	"sync"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/mullvad/message-queue/pubsub"
)

// This tests assumes that there's a sentinel running locally on 127.0.0.1:26379, with a group named "group"
// Both the sentinels and the redis servers should have authentication enabled, with the password "foobar"

const (
	redisAddress  = "redis://127.0.0.1:6379"
	redisPassword = "p4ssw0rd"
	channel       = "test"
	message       = "foobar"
)

func TestPubSub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}

	p, err := pubsub.New(redisAddress, redisPassword)
	if err != nil {
		t.Fatal(err)
	}

	defer p.Shutdown()

	ch, err := p.Subscribe(channel)
	if err != nil {
		t.Fatal(err)
	}
	assertReceiveMessages(t, ch)
}

func assertReceiveMessages(t *testing.T, ch <-chan []byte) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		actual := <-ch
		if string(actual) != message {
			t.Error("invalid message")
		}
		wg.Done()
	}()
	sendMessage(t)
	wg.Wait()
}

func sendMessage(t *testing.T) {
	t.Helper()

	conn, err := radix.Dial("tcp", redisAddress, radix.DialAuthPass(redisPassword))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.Do(radix.Cmd(nil, "PUBLISH", channel, message))
}
