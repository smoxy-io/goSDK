package logs

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestDefaultLogLevel(t *testing.T) {
	// catch accidental constant modification
	if DefaultLogLevel != "warn" {
		t.Errorf("DefaultLogLevel = %v, wanted: warn", DefaultLogLevel)
	}
}

func TestDefaultEncoding(t *testing.T) {
	// catch accidental constant modification
	if Default != JSON {
		t.Errorf("DefaultEncoding = %v, wanted: %v", Default, JSON)
	}
}

func TestNewLevel(t *testing.T) {
	validLevels := []string{
		"error",
		"warn",
		"info",
		"debug",
	}

	invalidLevels := []string{
		"foo",
		"bar",
		"baz",
		"warning",
		"debugging",
	}

	runTestNewLevel(validLevels, true, t)
	runTestNewLevel(invalidLevels, false, t)
}

func runTestNewLevel(tests []string, succeed bool, t *testing.T) {
	for _, s := range tests {
		l, err := NewLevel(s)

		if succeed {
			// expecting a valid level
			if err != nil {
				t.Errorf("NewLevel(%v) returned error: %v", s, err)
			}

			if l.String() != s {
				t.Errorf("NewLevel(%v) = %v, wanted: %v", s, l.String(), s)
			}
		} else {
			// expecting an invalid level
			if err == nil {
				t.Errorf("NewLevel(%v) did not return an error", s)
			}
		}
	}
}

func TestNewEncoderConfig(t *testing.T) {
	validateTestNewEncoderConfig(NewEncoderConfig(Console), Console, t)
	validateTestNewEncoderConfig(NewEncoderConfig(JSON), JSON, t)
	validateTestNewEncoderConfig(NewEncoderConfig(Encoding("foo")), JSON, t)
}

func validateTestNewEncoderConfig(c zapcore.EncoderConfig, e Encoding, t *testing.T) {
	if c.LevelKey != "level" {
		t.Errorf("LevelKey = %v, wanted: %v", c.LevelKey, "level")
	}

	if c.MessageKey != "message" {
		t.Errorf("MessageKey = %v, wanted: %v", c.MessageKey, "message")
	}

	if e == Console {
		// covered all settings for console encoding
		return
	}

	// remaining settings should always be checked for JSON encoding
	if c.TimeKey != "time" {
		t.Errorf("TimeKey = %v, wanted: %v", c.TimeKey, "time")
	}

	if c.CallerKey != "caller" {
		t.Errorf("CallerKey = %v, wanted: %v", c.CallerKey, "caller")
	}
}
