package EventBus

import (
	"github.com/smoxy-io/goSDK/util/events"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	New()

	if eventRouter == nil {
		t.Errorf("eventRouter = %v, wanted: *EventRouter", eventRouter)
	}

	// reset
	eventRouter = nil
}

func TestStop(t *testing.T) {
	New()
	Stop()

	// reset
	eventRouter = nil
}

func TestPublish(t *testing.T) {
	// start event bus
	New()

	if err := Publish("foo.bar", "test"); err != nil {
		t.Errorf("Publish('foo.bar', 'test') = '%v', wanted: '%v'", err, nil)
	}

	Stop()

	err := Publish("foo.bar", "test")

	if err == nil || err.Error() != "cannot publish event.  event router not started" {
		t.Errorf("error = '%v', wanted = '%v'", err, "cannot publish event.  event router not started")
	}

	err = Publish("foo.*.bar", "test")

	if err == nil || !strings.Contains(err.Error(), "invalid routing key") {
		t.Errorf("error = '%v', wanted = '%v'", err, "invalid routing key")
	}

	// reset
	eventRouter = nil
}

func TestSubscribe(t *testing.T) {
	// start event bus
	New()

	sub, err := Subscribe("foo.*")

	if err != nil {
		t.Errorf("error subscribings to 'foo.*', error: '%v'", err)
	}

	if sub == nil {
		t.Errorf("subscription not returned, got 'nil'")
	}

	Stop()

	_, err = Subscribe("foo.*")

	if err == nil || err.Error() != "event router not started" {
		t.Errorf("error = '%v', wanted = '%v'", err, "event router not started")
	}

	_, err = Subscribe("foo.*.*")

	if err == nil || err.Error() != "invalid topic" {
		t.Errorf("error = '%v', wanted = '%v'", err, "invalid topic")
	}

	// reset
	eventRouter = nil
}

func TestUnSubscribe(t *testing.T) {
	// start event bus
	New()

	sub, _ := Subscribe("foo.*")

	err := Unsubscribe("foo.*", sub)

	if err != nil {
		t.Errorf("error unsubscribings from 'foo.*', error: '%v'", err)
	}

	Stop()

	err = Unsubscribe("foo.*.*", sub)

	if err == nil || err.Error() != "invalid topic" {
		t.Errorf("error = '%v', wanted = '%v'", err, "invalid topic")
	}

	err = Unsubscribe("foo.*", sub)

	if err == nil || err.Error() != "event router not started" {
		t.Errorf("error = '%v', wanted = '%v'", err, "event router not started")
	}

	// reset
	eventRouter = nil
}

func TestSendReceive(t *testing.T) {
	New()

	topic1 := "*"
	topic2 := "foo.*"
	routingKey := "foo.bar"
	event := "test1"

	subsReady := sync.WaitGroup{}
	wg := sync.WaitGroup{}

	subFunc := func(name, topic string) {
		defer wg.Done()

		sub, err := Subscribe(topic)

		if err != nil {
			t.Errorf("%v: failed to subscribe to '%v', error: %v", name, topic, err)
			subsReady.Done()
			return
		}

		subsReady.Done()

		ev, open := <-sub

		if !open {
			t.Errorf("%v: subscription closed before receiving event", name)
		}

		if ok, err := ev.IsValid(); !ok || err != nil || !reflect.DeepEqual(UnwrapEvent[string](ev), event) {
			t.Errorf("%v: received incorrect event: %v, wanted: %v, error: %v", name, ev, event, err)
		}

		err = Unsubscribe(topic, sub)

		if err != nil {
			t.Errorf("%v: error unsubscribing from topic '%v', error: %v", name, topic, err)
		}
	}

	// subscriber 1
	wg.Add(1)
	subsReady.Add(1)
	go subFunc("sub1", topic1)

	// subscriber 2
	wg.Add(1)
	subsReady.Add(1)
	go subFunc("sub2", topic1)

	// subscriber 3
	wg.Add(1)
	subsReady.Add(1)
	go subFunc("sub3", topic2)

	subsReady.Wait()

	// publisher 1
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := Publish(routingKey, event); err != nil {
			t.Errorf("publisher: got error publishing event: '%v', error: %v", event, err)
		}
	}()

	wg.Wait()

	Stop()

	// reset
	eventRouter = nil
}

func TestUnwrapEvent(t *testing.T) {
	if UnwrapEvent[string](events.NewEvent(events.RoutingKey("foo"), "bar")) != "bar" {
		t.Errorf("expected: %v, got: %v", "bar", UnwrapEvent[string](events.NewEvent(events.RoutingKey("foo"), "bar")))
	}

	if UnwrapEvent[int](events.NewEvent(events.RoutingKey("foo"), 1)) != 1 {
		t.Errorf("expected: %v, got: %v", 1, UnwrapEvent[int](events.NewEvent(events.RoutingKey("foo"), 1)))
	}
}
