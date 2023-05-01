package events

import "regexp"

type RoutingKey string

const (
	ValidRoutingKeyPattern = `^\w[-\w]*(\.\w[-\w]*)*$`
)

func (rk RoutingKey) IsValid() bool {
	valid, err := regexp.Match(ValidRoutingKeyPattern, []byte(string(rk)))

	if err != nil {
		return false
	}

	return valid
}

func (rk RoutingKey) String() string {
	return string(rk)
}
