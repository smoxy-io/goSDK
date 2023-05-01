package events

import (
	"testing"
	"time"
)

func TestEvent_IsValid(t *testing.T) {
	validEvents := []Event{
		// all fields are valid
		{
			RoutingKey: "foo",
			Msg:        1,
			Timestamp:  time.Now(),
		},
	}

	invalidEvents := []Event{
		// all fields are invalid
		{
			RoutingKey: "",
			Msg:        nil,
			Timestamp:  time.Now().Add(time.Hour * 100), // future time
		},
	}

	runTestEvent_IsValid(validEvents, true, t)
	runTestEvent_IsValid(invalidEvents, false, t)
}

func runTestEvent_IsValid(tests []Event, expected bool, t *testing.T) {
	for _, test := range tests {
		if res, err := test.IsValid(); res != expected {
			t.Errorf("Event.IsValid() = %v, wanted %v (event: %v, err: %v)", res, expected, test, err)

			// all invalid results MUST have an accompanying error
			if !res && err == nil {
				t.Errorf("when Event.IsValid() returns 'false' it MUST also return an error")
			}
		}
	}
}
