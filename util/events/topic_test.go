package events

import (
	"regexp"
	"testing"
)

func TestTopic_IsValid(t *testing.T) {
	validTopics := []Topic{
		"a",
		"_",
		"1",
		"a.b.c",
		"a1_3-2.foo.blah.12-",
		"*",
		"foo.#",
		"foo.#.bar",
		"foo.#.var.*",
		"baz.#.foo.#.bar",
	}

	invalidTopics := []Topic{
		".",
		"-",
		"-.#",
		"abc..def",
		"123.-",
		"a.#*",
		"a.*#",
		"#.foo.#.bar",
		"a.**",
		"b.*.*",
		"b.#.b.#.*.*",
		"*.a",
	}

	runTestTopic_IsValid(validTopics, true, t)
	runTestTopic_IsValid(invalidTopics, false, t)
}

func runTestTopic_IsValid(tests []Topic, expected bool, t *testing.T) {
	for _, test := range tests {
		if v := test.IsValid(); v != expected {
			t.Errorf("'%v'.IsValid() = %v, wanted %v", test, v, expected)
		}
	}
}

func TestTopic_Matches(t *testing.T) {
	shouldMatch := []map[Topic][]RoutingKey{
		{
			Topic("*"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("lorim.ipsum"),
				RoutingKey("lorim-ipsum"),
				RoutingKey("lorim-ipsum.foo.bar-baz"),
			},
			Topic("foo.#"): {
				RoutingKey("foo.b"),
				RoutingKey("foo.c"),
				RoutingKey("foo.cd"),
				RoutingKey("foo.cdefg"),
				RoutingKey("foo.bar"),
				RoutingKey("foo.bar-baz"),
			},
			Topic("foo.*"): {
				RoutingKey("foo.b"),
				RoutingKey("foo.b.c"),
				RoutingKey("foo.b.c.d"),
				RoutingKey("foo.cd"),
				RoutingKey("foo.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("foo.lorim.ipsum"),
				RoutingKey("foo.lorim-ipsum"),
				RoutingKey("foo.lorim-ipsum.foo.bar-baz"),
			},
			Topic("foo.#.bar"): {
				RoutingKey("foo.cd.bar"),
				RoutingKey("foo.bar.bar"),
				RoutingKey("foo.lorim-ipsum.bar"),
			},
			Topic("foo.#.bar.*"): {
				RoutingKey("foo.b.bar.a.b"),
				RoutingKey("foo.a.bar.c"),
				RoutingKey("foo.c.bar.c.d"),
				RoutingKey("foo.baz.bar.bar"),
				RoutingKey("foo.lorim.bar.ipsum"),
				RoutingKey("foo.lorim-ipsum.bar.baz"),
			},
			Topic("foo.#.bar.#"): {
				RoutingKey("foo.b.bar.a"),
				RoutingKey("foo.a.bar.c"),
				RoutingKey("foo.c.bar.d"),
				RoutingKey("foo.bar.bar.baz"),
				RoutingKey("foo.lorim.bar.ipsum"),
				RoutingKey("foo.lorim-ipsum.bar.baz"),
			},
			Topic("foo.#.bar.#.baz"): {
				RoutingKey("foo.b.bar.a.baz"),
				RoutingKey("foo.a.bar.c.baz"),
				RoutingKey("foo.c.bar.d.baz"),
				RoutingKey("foo.bar.bar.baz.baz"),
				RoutingKey("foo.lorim.bar.ipsum.baz"),
				RoutingKey("foo.lorim-ipsum.bar.baz.baz"),
			},
		},
	}

	shouldNotMatch := []map[Topic][]RoutingKey{
		{
			Topic("foo.#"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("foo"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("foo.lorim.ipsum"),
				RoutingKey("foo.lorim-ipsum.foo.bar-baz"),
			},
			Topic("foo.*"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("lorim.ipsum"),
				RoutingKey("lorim-ipsum"),
				RoutingKey("lorim-ipsum.foo.bar-baz"),
			},
			Topic("foo.#.bar"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("lorim.ipsum"),
				RoutingKey("lorim-ipsum"),
				RoutingKey("lorim-ipsum.foo.bar-baz"),
			},
			Topic("foo.#.bar.*"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("lorim.ipsum"),
				RoutingKey("lorim-ipsum"),
				RoutingKey("lorim-ipsum.foo.bar-baz"),
				RoutingKey("foo.baz.lorim.bar"),
			},
			Topic("foo.#.bar.#"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("lorim.ipsum"),
				RoutingKey("lorim-ipsum"),
				RoutingKey("lorim-ipsum.foo.bar-baz"),
				RoutingKey("foo.baz.bar.lorim.ipsum"),
			},
			Topic("foo.#.bar.#.baz"): {
				RoutingKey("a"),
				RoutingKey("a.b"),
				RoutingKey("a.b.c"),
				RoutingKey("a.b.c.d"),
				RoutingKey("ab.cd"),
				RoutingKey("ab.cd.efg"),
				RoutingKey("foo.bar.baz"),
				RoutingKey("lorim.ipsum"),
				RoutingKey("lorim-ipsum"),
				RoutingKey("lorim-ipsum.foo.bar-baz"),
				RoutingKey("foo.baz.bar.lorim.baz.ipsum"),
			},
		},
	}

	runTestTopic_Matches(shouldMatch, true, t)
	runTestTopic_Matches(shouldNotMatch, false, t)
}

func runTestTopic_Matches(tests []map[Topic][]RoutingKey, expected bool, t *testing.T) {
	for _, topics := range tests {
		for topic, routingKeys := range topics {
			for _, routingKey := range routingKeys {
				if match := topic.Matches(routingKey); match != expected {
					t.Errorf("'%v'.Matches('%v') = %v, wanted %v", topic, routingKey, match, expected)
				}
			}
		}
	}
}

func TestTopic_ToRegexp(t *testing.T) {
	tests := map[Topic]*regexp.Regexp{
		Topic("*"):       regexp.MustCompile(`^.*$`),
		Topic("a.*"):     regexp.MustCompile(`^a\..*$`),
		Topic("a.#"):     regexp.MustCompile(`^a\.[^.]+$`),
		Topic("a.#.b.*"): regexp.MustCompile(`^a\.[^.]+\.b\..*$`),
		Topic("a.#.b.#"): regexp.MustCompile(`^a\.[^.]+\.b\.[^.]+$`),
	}

	for topic, reg := range tests {
		if r := topic.ToRegexp(); r.String() != reg.String() {
			t.Errorf("'%v'.ToRegexp() = /%v/, wanted: /%v/", topic, r, reg)
		}
	}
}

func TestTopic_String(t *testing.T) {
	tests := map[Topic]string{
		Topic("foo"):             "foo",
		Topic("foo.*"):           "foo.*",
		Topic("foo.#"):           "foo.#",
		Topic("foo.#.bar"):       "foo.#.bar",
		Topic("foo.#.bar-baz.*"): "foo.#.bar-baz.*",
	}

	for topic, str := range tests {
		if tstr := topic.String(); tstr != str {
			t.Errorf("Topic.String() = '%v', wanted: '%v'", tstr, str)
		}
	}
}
