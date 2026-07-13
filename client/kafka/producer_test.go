package kafkaclient

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestProducerRoundWaitsForMessageOrCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Producer{
		KafkaClient: &KafkaClient{ctx: ctx},
		chSend:      make(chan *Message, 1),
	}

	result := make(chan error, 1)
	go func() {
		result <- p.round()
	}()

	select {
	case err := <-result:
		t.Fatalf("round returned while the send queue was empty: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("round returned %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("round did not return after context cancellation")
	}
}

func TestProducerProduceReturnsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Producer{
		KafkaClient: &KafkaClient{ctx: ctx},
		chSend:      make(chan *Message, 1),
	}

	if err := p.Produce(&Message{}); err != nil {
		t.Fatalf("first Produce returned an unexpected error: %v", err)
	}

	result := make(chan error, 1)
	go func() {
		result <- p.Produce(&Message{})
	}()

	select {
	case err := <-result:
		t.Fatalf("Produce returned while the send queue was full: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Produce returned %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Produce did not return after context cancellation")
	}
}

func TestProducerProduceDoesNotEnqueueAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Producer{
		KafkaClient: &KafkaClient{ctx: ctx},
		chSend:      make(chan *Message, 1),
	}
	cancel()

	if err := p.Produce(&Message{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("Produce returned %v, want context.Canceled", err)
	}
	if got := len(p.chSend); got != 0 {
		t.Fatalf("send queue contains %d messages after cancellation, want 0", got)
	}
}
