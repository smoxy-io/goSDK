package Logger

import (
	"github.com/smoxy-io/goSDK/util/logs"
	"go.uber.org/zap"
	"testing"
)

func TestInitLogger(t *testing.T) {
	log, err := InitLogger("io", "info", "console")

	if err != nil {
		t.Errorf("InitLogger() error = %q, wanted %v", err, nil)
	}

	if ent := log.Check(zap.InfoLevel, "foo bar"); ent == nil {
		t.Error("Default loglevel should allow 'info' logs")
	}

	if ent := log.Check(zap.WarnLevel, "foo bar"); ent == nil {
		t.Error("Default loglevel should allow 'warn' logs")
	}

	if ent := log.Check(zap.ErrorLevel, "foo bar"); ent == nil {
		t.Error("Default loglevel should allow 'error' logs")
	}

	if ent := log.Check(zap.DebugLevel, "foo bar"); ent != nil {
		t.Error("Default loglevel should not allow 'debug' logs")
	}
}

func TestGetLogger(t *testing.T) {
	l := GetLogger()

	if l == nil {
		t.Error("Expected *zap.Logger got nil")
	}

	if l != logger {
		t.Errorf("Wanted %v, got %v", logger, l)
	}
}

func TestRegisterLoggerInit(t *testing.T) {
	RegisterLoggerInit("foo", InitIOLogger)

	init, ok := loggers["foo"]

	if !ok {
		t.Error("failed to register logger init function")
	}

	if init == nil {
		t.Error("no logger init function found at key 'foo'")
	}
}

func TestInitIOLogger(t *testing.T) {
	var err error

	_, err = InitIOLogger(logs.DefaultLogLevel, logs.JSON)

	if err != nil {
		t.Error(err)
	}

	_, err = InitIOLogger(logs.DefaultLogLevel, logs.Console)

	if err != nil {
		t.Error(err)
	}

	_, err = InitIOLogger("foo", logs.JSON)

	if err == nil {
		t.Errorf("did not receive error with invalid log level: %v", "foo")
	}

	_, err = InitIOLogger(logs.DefaultLogLevel, "bar")

	if err != nil {
		t.Errorf("expected invalid encoding to be silently corrected to JSON encoding.  instead got error: %v", err)
	}
}
