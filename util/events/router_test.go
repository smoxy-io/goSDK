package events

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewEventRouter(t *testing.T) {
	_ = NewEventRouter()
}

func TestEventRouter_Start_Stop(t *testing.T) {
	er := NewEventRouter()

	er.Start()
	er.Stop()
}

func TestEventRouter_Subscribe_Publish_Unsubscribe(t *testing.T) {
	er := NewEventRouter()

	er.Start()

	topic := Topic("*")
	subscription, err := er.Subscribe(topic)

	if err != nil {
		t.Errorf("error subscribing to topic: '%v'.  error: '%v'", topic, err)
		er.Stop()
		return
	}

	if subs, ok := er.subscribers[topic]; !ok || len(subs) != 1 || subs[0].Subscriber != subscription {
		t.Errorf("subscriptions not being tracked correctly after subscribe.  EventRouter.subscribers: %v", er.subscribers)
	}

	routingKey := RoutingKey("foo.bar")
	msg := 1
	if err := er.Publish(routingKey, msg); err != nil {
		t.Errorf("error publishing event '%v' with routingKey: '%v'.  error: '%v'", msg, routingKey, err)
		er.Stop()
		return
	}

	ev := <-subscription

	if err := er.Unsubscribe(topic, subscription); err != nil {
		t.Errorf("error unsubscribing from topic: '%v'.  error: '%v'", topic, err)
		er.Stop()
		return
	}

	if len(er.subscribers) != 0 {
		t.Errorf("subscriptions not being tracked correctly after unsubscribe.  EventRouter.subscribers: %v", er.subscribers)
	}

	er.Stop()

	if ok, err := ev.IsValid(); !ok {
		t.Errorf("invalid event received: %v, error: %v", ev, err)
	}

	if ev.Msg != msg {
		t.Errorf("subscriber received '%v', wanted: '%v'", ev.Msg, msg)
	}

	if ev.RoutingKey != routingKey {
		t.Errorf("subscriber received different routing key: '%v'.  wanted: '%v'", ev.RoutingKey, routingKey)
	}

	if !topic.Matches(ev.RoutingKey) {
		t.Errorf("subscriber received event it was not subscribed to.  topic: '%v', routingKey: '%v'", topic, ev.RoutingKey)
	}
}

func TestEventRouter_SubscribeMultiple(t *testing.T) {
	var resAll1, resAll2, resFoo1, resFoo2 []Event

	topic1 := Topic("*")
	topic2 := Topic("foo.#")

	events := map[Topic][]Event{
		topic1: {
			NewEvent(RoutingKey("foo"), "bar"),
			NewEvent(RoutingKey("foo.bar.baz"), "baz"),
			NewEvent(RoutingKey("lorim"), "ipsum"),
		},
		topic2: {
			NewEvent(RoutingKey("foo.bar"), "baz"),
			NewEvent(RoutingKey("foo.lorim"), "ipsum"),
		},
	}

	expectedTopic1 := len(events[topic1]) + len(events[topic2])
	expectedTopic2 := len(events[topic2])

	er := NewEventRouter()
	er.Start()

	subAll1, _ := er.Subscribe(topic1)
	subAll2, _ := er.Subscribe(topic1)

	if len(er.subscribers[topic1]) != 2 {
		t.Errorf("topic '%v' has %v subscriber(s), wanted %v", topic1, len(er.subscribers[topic1]), 2)
	}

	subFoo1, _ := er.Subscribe(topic2)
	subFoo2, _ := er.Subscribe(topic2)

	if len(er.subscribers[topic2]) != 2 {
		t.Errorf("topic '%v' has %v subscriber(s), wanted %v", topic2, len(er.subscribers[topic2]), 2)
	}

	if len(er.subscribers) != 2 {
		t.Errorf("%v topic(s) subscribed to, wanted %v", len(er.subscribers), 2)
	}

	st1, ok := (er.subscribers)[topic1]

	if !ok {
		t.Errorf("no subscribers for topic: %v", topic1)
	}

	if len(st1) != 2 {
		t.Errorf("len(EventRouter_subscribers[%v]) = %v, wanted %v", len(st1), topic1, 2)
	}

	wg := sync.WaitGroup{}

	consume := func(name string, subscription Subscriber, expected int, res *[]Event, t *testing.T) {
		defer wg.Done()

		count := 0

		for {
			event, open := <-subscription

			if !open {
				if count != expected {
					t.Errorf("subscription '%v' got %v events, wanted %v", name, count, expected)
				}

				return
			}

			// for debugging test
			//t.Logf("%v: %v", name, event)

			*res = append(*res, event)
			count++
		}
	}

	wg.Add(1)
	go consume("All1", subAll1, expectedTopic1, &resAll1, t)
	wg.Add(1)
	go consume("All2", subAll2, expectedTopic1, &resAll2, t)
	wg.Add(1)
	go consume("Foo1", subFoo1, expectedTopic2, &resFoo1, t)
	wg.Add(1)
	go consume("Foo2", subFoo2, expectedTopic2, &resFoo2, t)

	for _, evs := range events {
		for _, e := range evs {
			if err := er.PublishEvent(e); err != nil {
				t.Errorf("error publishing event: %v, error: %v", e, err)
			}
		}
	}

	er.Stop()

	if len(er.subscribers) != 0 {
		t.Errorf("%v topic(s) subscribed to, wanted %v", len(er.subscribers), 0)
	}

	wg.Wait()

	if len(resAll1) != expectedTopic1 {
		t.Errorf("subAll1 received %v events, wanted %v, events: %v", len(resAll1), expectedTopic1, resAll1)
	}

	if len(resAll2) != expectedTopic1 {
		t.Errorf("subAll2 received %v events, wanted %v, events: %v", len(resAll2), expectedTopic1, resAll2)
	}

	if len(resFoo1) != expectedTopic2 {
		t.Errorf("subFoo1 received %v events, wanted %v, events: %v", len(resFoo1), expectedTopic2, resFoo1)
	}

	if len(resFoo2) != expectedTopic2 {
		t.Errorf("subFoo2 received %v events, wanted %v, events: %v", len(resFoo2), expectedTopic2, resFoo2)
	}

	if !reflect.DeepEqual(resAll1, resAll2) {
		t.Errorf("All subscribers to topic '%v' should receive the same events.  %v != %v", topic1, resAll1, resAll2)
	}

	if !reflect.DeepEqual(resFoo1, resFoo2) {
		t.Errorf("All subscribers to topic '%v' should receive the same events.  %v != %v", topic2, resFoo1, resFoo2)
	}
}
