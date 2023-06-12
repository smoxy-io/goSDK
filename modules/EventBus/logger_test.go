package EventBus

import (
	"bytes"
	"fmt"
	"go.uber.org/zap/zapcore"
	"strings"
	"testing"
)

func TestInitLogger(t *testing.T) {
	New()

	logger, err := InitLogger("info", "json")

	if err != nil {
		t.Errorf("InitLogger() returned an error when no error was expected.  error: %v", err)
	}

	if !logger.Level().Enabled(zapcore.InfoLevel) {
		t.Errorf("'info' level logs are not enabled")
	}

	logger.Info("event bus logger")

	Stop()

	// reset
	eventRouter = nil
}

func TestGenerateRoutingKey(t *testing.T) {
	e := zapcore.Entry{}
	wanted := fmt.Sprintf("%v%v", RoutingKeyBase, e.Level.String())

	rk := generateRoutingKey(e)

	if rk != wanted {
		t.Errorf("wanted: %v, got: %v", wanted, rk)
	}

	e.LoggerName = "myApp"
	wantedWithName := fmt.Sprintf("%v.app.%v", wanted, strings.Replace(e.LoggerName, ".", "-", -1))

	rk = generateRoutingKey(e)

	if rk != wantedWithName {
		t.Errorf("wanted: %v, got: %v", wantedWithName, rk)
	}

	e.LoggerName = "my.App"
	wantedWithName = fmt.Sprintf("%v.app.%v", wanted, strings.Replace(e.LoggerName, ".", "-", -1))

	rk = generateRoutingKey(e)

	if rk != wantedWithName {
		t.Errorf("wanted: %v, got: %v", wantedWithName, rk)
	}
}

func TestIOLogSubscriber(t *testing.T) {
	New()

	buff := new(bytes.Buffer)

	subscriber, err := IOLogSubscriber(buff)

	if err != nil {
		t.Errorf("IOLogSubscriber() returned error: %v", err)
	}

	logger, err := InitLogger("info", "json")

	if err != nil {
		t.Errorf("InitLogger() returned an error when no error was expected.  error: %v", err)
	}

	// default logLevel is warn
	logger.Info("test")

	var data []byte

	for {
		data = buff.Bytes()

		if len(data) > 0 {
			break
		}
	}

	if len(data) != 72 {
		t.Errorf("wanted: %v, got: %v", 72, len(data))
	}

	if !strings.Contains(string(data), "\"level\":\"info\"") {
		t.Errorf("wanted log to contain: %v, got: %v", "\"level\":\"info\"", string(data))
	}

	if !strings.Contains(string(data), "\"message\":\"test\"") {
		t.Errorf("wanted log to contain: %v, got: %v", "\"message\":\"test\"", string(data))
	}

	if err := Unsubscribe(AllLogsTopic, subscriber); err != nil {
		t.Errorf("Unsubscribing the IOLoggerSubscription returned an error: %v", err)
	}

	Stop()

	// reset
	eventRouter = nil
}
