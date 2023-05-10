package events

import (
	"errors"
	"github.com/smoxy-io/goSDK/util/maps"
	"sync"
	"sync/atomic"
)

const (
	EventBufferSize      = 128
	SubscriberBufferSize = 64
)

type RoutingPair struct {
	Channel chan Event
	Publisher
	Subscriber
}

type EventRouter struct {
	subscribers     map[Topic][]RoutingPair
	topicMatchCache map[RoutingKey][]*Topic
	eventChan       chan Event
	eventWg         *sync.WaitGroup
	_subscribers    atomic.Pointer[map[Topic][]RoutingPair]
}

func NewEventRouter() *EventRouter {
	evr := EventRouter{
		eventWg:     new(sync.WaitGroup),
		subscribers: map[Topic][]RoutingPair{},
	}

	// initialize the _subscribers member
	evr.reload()

	return &evr
}

func (er *EventRouter) Start() {
	if er.eventChan != nil {
		return
	}

	er.eventChan = make(chan Event, EventBufferSize)
	er.eventWg.Add(1)

	go er.routeEvents(er.eventChan)

	er.eventWg.Wait()
}

func (er *EventRouter) Stop() {
	if er.eventChan == nil {
		return
	}

	er.eventWg.Add(1)

	close(er.eventChan)

	er.eventWg.Wait()
	er.eventChan = nil

	// clean up subscribers
	er.unsubscribeAll()
}

func (er *EventRouter) Subscribe(topic Topic) (Subscriber, error) {
	if !topic.IsValid() {
		return nil, errors.New("invalid topic")
	}

	if er.eventChan == nil {
		return nil, errors.New("event router not started")
	}

	subscription := make(chan Event, SubscriberBufferSize)
	er.subscribers[topic] = append(er.subscribers[topic], RoutingPair{
		Channel:    subscription,
		Publisher:  subscription,
		Subscriber: subscription,
	})

	er.reload()

	return subscription, nil
}

func (er *EventRouter) Unsubscribe(topic Topic, subscription Subscriber) error {
	if !topic.IsValid() {
		return errors.New("invalid topic")
	}

	if er.eventChan == nil {
		return errors.New("event router not started")
	}

	err := er.removeSubscriber(topic, subscription)

	if err != nil {
		return err
	}

	er.reload()

	return nil
}

func (er *EventRouter) removeSubscriber(topic Topic, subscription Subscriber) error {
	subscriptions, ok := er.subscribers[topic]

	if !ok {
		// no subscribers for this topic.  no need to reload handlers
		return nil
	}

	sIndex := -1
	var pair RoutingPair

	for i, p := range subscriptions {
		if p.Subscriber == subscription {
			// found the subscriber to remove
			sIndex = i
			pair = p
			break
		}
	}

	if sIndex < 0 {
		// no matching handler to remove.  no need to reload handlers
		return nil
	}

	// close the channel
	close(pair.Channel)

	er.subscribers[topic] = append(er.subscribers[topic][:sIndex], er.subscribers[topic][sIndex+1:]...)

	if len(er.subscribers[topic]) < 1 {
		// remove the topic if there are no more subscribers
		delete(er.subscribers, topic)
	}

	return nil
}

func (er *EventRouter) reload() {
	tmpSubscribers := maps.Clone(er.subscribers)
	er._subscribers.Swap(&tmpSubscribers)
}

func (er *EventRouter) unsubscribeAll() {
	for t, pairs := range er.subscribers {
		for _, p := range pairs {
			// ignore unsubscribe errors
			_ = er.removeSubscriber(t, p.Subscriber)
		}
	}

	er.reload()
}

// this is the main event loop.  it is run inside a go routine
func (er *EventRouter) routeEvents(eventChan <-chan Event) {
	// let the stopping function know when we've finished
	defer er.eventWg.Done()

	// let the starting function know we've started
	er.eventWg.Done()

	for {
		select {
		case event, open := <-eventChan:
			if !open {
				// channel closed.  exit
				return
			}

			// check and wait for an active subscriber reload
			er.routeEvent(event)
		}
	}
}

func (er *EventRouter) routeEvent(event Event) {
	// load the current subscriber list
	subscribers := *er._subscribers.Load()

	if len(subscribers) < 1 {
		// no subscribers
		return
	}

	// get the topics
	topics := maps.Keys(subscribers)

	wg := sync.WaitGroup{}

	for _, t := range topics {
		if t.Matches(event.RoutingKey) {
			// send the event to all subscribers of the matching topic
			// process each topic's subscribers in its own go routine
			wg.Add(1)
			go func(waitgroup *sync.WaitGroup, subs []RoutingPair) {
				defer waitgroup.Done()
				for _, pair := range subs {
					pair.Publisher <- event
				}
			}(&wg, subscribers[t])
		}
	}

	wg.Wait()
}

func (er *EventRouter) Publish(routingKey RoutingKey, event any) error {
	return er.PublishEvent(NewEvent(routingKey, event))
}

func (er *EventRouter) PublishEvent(event Event) error {
	if ok, err := event.IsValid(); !ok {
		// invalid event
		return err
	}

	if er.eventChan == nil {
		// can't publish messages when the routing is not running
		return errors.New("cannot publish event.  event router not started")
	}

	er.eventChan <- event

	return nil
}
