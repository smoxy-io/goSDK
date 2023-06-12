package EventBus

import (
	"fmt"
	"github.com/smoxy-io/goSDK/util/arrays"
	"github.com/smoxy-io/goSDK/util/events"
	"github.com/smoxy-io/goSDK/util/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"sync"
)

const RoutingKeyBase = "type.log.level."
const AllLogsTopic = RoutingKeyBase + "*"
const DebugLogsTopic = RoutingKeyBase + "debug.*"
const InfoLogsTopic = RoutingKeyBase + "info.*"
const WarnLogsTopic = RoutingKeyBase + "warn.*"
const ErrorLogsTopic = RoutingKeyBase + "error.*"

type LogEventBuffer []byte

type loggerCore struct {
	level    zapcore.LevelEnabler
	encoder  zapcore.Encoder
	out      zapcore.WriteSyncer
	encoding string
}

func (e *loggerCore) Sync() error {
	return e.out.Sync()
}

func (e *loggerCore) Enabled(level zapcore.Level) bool {
	return e.level.Enabled(level)
}

func (e *loggerCore) With(fields []zapcore.Field) zapcore.Core {
	c := e.clone()
	addFields(e.encoder, fields)
	return c
}

func (e *loggerCore) Check(entry zapcore.Entry, entry2 *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if e.Enabled(entry.Level) {
		return entry2.AddCore(entry, e)
	}

	return entry2
}

func (e *loggerCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buff, err := e.encoder.EncodeEntry(entry, fields)

	if err != nil {
		return err
	}

	rk := generateRoutingKey(entry)

	err = Publish(rk, LogEventBuffer(arrays.Clone(buff.Bytes())))
	buff.Free()

	if err != nil {
		return err
	}

	if entry.Level > zap.ErrorLevel {
		_ = e.Sync()
	}

	return nil
}

func generateRoutingKey(entry zapcore.Entry) string {
	rk := fmt.Sprintf("%v%v", RoutingKeyBase, entry.Level.String())

	if entry.LoggerName != "" {
		rk = fmt.Sprintf("%v.app.%v", rk, strings.Replace(entry.LoggerName, ".", "-", -1))
	}

	return rk
}

func (e *loggerCore) clone() *loggerCore {
	return &loggerCore{
		level:   e.level,
		encoder: e.encoder.Clone(),
		out:     zapcore.AddSync(io.Discard),
	}
}

func InitLogger(logLevel string, encoding string, options ...zap.Option) (*zap.Logger, error) {
	var newEncoder func(cfg zapcore.EncoderConfig) zapcore.Encoder

	cfg, err := logs.NewLoggerConfig(logLevel, logs.Encoding(encoding))

	if err != nil {
		return nil, err
	}

	if logs.Encoding(cfg.Encoding) == logs.Console {
		newEncoder = zapcore.NewConsoleEncoder
	} else {
		newEncoder = zapcore.NewJSONEncoder
	}

	logger := zap.New(newLoggerCore(newEncoder(cfg.EncoderConfig), cfg.Level.Level(), cfg.Encoding), options...)

	return logger, nil
}

func newLoggerCore(encoder zapcore.Encoder, level zapcore.Level, encoding string) zapcore.Core {
	return &loggerCore{
		level:    level,
		encoder:  encoder,
		out:      zapcore.AddSync(io.Discard),
		encoding: encoding,
	}
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

// IOLogSubscriber creates an eventbus subscriber that writes all logs from the eventbus to the IO Writer w.
// the events.Subscriber returned from this function is for advanced usage and can be safely ignored in 99% of
// use cases.
//
// Example:
//
//	_, err := IOLogSubscriber(ioWriter)
//	if err != nil {
//	  // failed to create the log subscriber
//	}
//
//	// continue with application functions
func IOLogSubscriber(w io.Writer) (events.Subscriber, error) {
	// subscribe to all logs
	return ioLogSubscriber(AllLogsTopic, w)
}

func ioLogSubscriber(topic string, w io.Writer) (events.Subscriber, error) {
	sub, err := Subscribe(topic)

	if err != nil {
		// can't subscribe
		return nil, err
	}

	subReady := &sync.WaitGroup{}
	subReady.Add(1)

	go ioLogProcessor(sub, w, subReady)

	subReady.Wait()

	return sub, nil
}

func ioLogProcessor(sub events.Subscriber, w io.Writer, ready *sync.WaitGroup) {
	ready.Done()

	for {
		e, ok := <-sub

		if !ok {
			// subscription channel is closed.  exit go routine
			return
		}

		if v, err := e.IsValid(); !v || err != nil {
			// invalid event received.  log this error instead
			_ = Publish(e.RoutingKey.String(), LogEventBuffer("log subscriber received invalid event: "+e.String()))
			continue
		}

		// write the log to the io channel
		ioWriter(e, w)
	}
}

func ioWriter(e events.Event, w io.Writer) {
	d := UnwrapEvent[LogEventBuffer](e)
	l := len(d)

	for {
		b, err := w.Write(d)

		if err != nil {
			// log the error
			_ = Publish(e.RoutingKey.String(), LogEventBuffer("error writing log to io channel.  data: "+string(d)+", bytes written: "+string(rune(b))+"/"+string(rune(l))+", error: "+err.Error()))
			break
		}

		if b == 0 {
			// nothing written, but no error
			_ = Publish(e.RoutingKey.String(), LogEventBuffer("no bytes written to io channel.  data: "+string(d)+", bytes written: "+string(rune(b))+"/"+string(rune(l))+", error: "+err.Error()))
			break
		}

		if b != l {
			d = d[b:]
			l = l - b
			continue
		}

		break
	}
}

//
// Convenience functions for common IO channels where logs are written
//

func SplitStdoutStderrLogSubscriber() ([]events.Subscriber, error) {
	var subs []events.Subscriber

	sub1, err1 := ioLogSubscriber(ErrorLogsTopic, os.Stderr)
	sub2, err2 := ioLogSubscriber(WarnLogsTopic, os.Stderr)
	sub3, err3 := ioLogSubscriber(InfoLogsTopic, os.Stdout)
	sub4, err4 := ioLogSubscriber(DebugLogsTopic, os.Stdout)

	if err1 != nil {
		return subs, err1
	}

	if err2 != nil {
		return subs, err2
	}

	if err3 != nil {
		return subs, err3
	}

	if err4 != nil {
		return subs, err4
	}

	subs = append(subs, sub1, sub2, sub3, sub4)

	return subs, nil
}

func StdoutLogSubscriber() (events.Subscriber, error) {
	return IOLogSubscriber(os.Stdout)
}

func StderrLogSubscriber() (events.Subscriber, error) {
	return IOLogSubscriber(os.Stderr)
}

func FileLogSubscriber(file *os.File) (events.Subscriber, error) {
	return IOLogSubscriber(file)
}
