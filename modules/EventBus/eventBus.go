package EventBus

import (
	"goSDK/util/events"
)

// event router for the event bus
var eventRouter *events.EventRouter

func New() {
	if eventRouter != nil {
		// event router already created.  nothing to do
		return
	}

	// create the event router for the event bus
	eventRouter = events.NewEventRouter()

	// start the event router
	eventRouter.Start()
}

func Publish(routingKey string, event any) error {
	return eventRouter.Publish(events.RoutingKey(routingKey), event)
}

func Subscribe(topic string) (events.Subscriber, error) {
	return eventRouter.Subscribe(events.Topic(topic))
}

func Unsubscribe(topic string, subscription events.Subscriber) error {
	return eventRouter.Unsubscribe(events.Topic(topic), subscription)
}

func UnwrapEvent[T any](event events.Event) T {
	return event.Msg.(T)
}

func Stop() {
	eventRouter.Stop()
}
