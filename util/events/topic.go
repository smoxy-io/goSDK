package events

import (
	"regexp"
	"strings"
)

type Topic string

const (
	TopicMatchAll  = "*"
	TopicMatchPart = "#"
)

const (
	ValidTopicPattern = `^(\*|\w[-\w]*(\.(#|\w[-\w]*))*(\.\*)?)$`
	TopicSeparator    = "."
)

func (t Topic) Matches(routingKey RoutingKey) bool {
	if !routingKey.IsValid() {
		// invalid routing key (never matches)
		return false
	}

	tStr := string(t)

	if tStr == TopicMatchAll {
		// matches all routing keys
		return true
	}

	rkStr := string(routingKey)

	if tStr == rkStr {
		// exact match to routing key
		return true
	}

	if !strings.Contains(tStr, TopicMatchPart) && !strings.Contains(tStr, TopicMatchAll) {
		// topic only matches an exact match
		return false
	}

	// create a regexp from the topic
	reg := t.ToRegexp()

	return reg.Match([]byte(rkStr))
}

func (t Topic) IsValid() bool {
	valid, err := regexp.Match(ValidTopicPattern, []byte(string(t)))

	if err != nil {
		return false
	}

	return valid
}

func (t Topic) ToRegexp() *regexp.Regexp {
	pattern := strings.ReplaceAll(string(t), TopicSeparator, `\`+TopicSeparator)
	pattern = strings.ReplaceAll(pattern, TopicMatchPart, `[^`+TopicSeparator+`]+`)

	r, _ := regexp.Compile(`[` + TopicMatchAll + `].*$`)

	pattern = "^" + string(r.ReplaceAll([]byte(pattern), []byte(`.*`))) + "$"

	reg, _ := regexp.Compile(pattern)

	return reg
}

func (t Topic) String() string {
	return string(t)
}
