package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerInit func(logLevel string, encoding Encoding, options ...zap.Option) (*zap.Logger, error)

type Encoding string

const (
	Console Encoding = "console"
	JSON    Encoding = "json"
	Default Encoding = JSON
)

const (
	DefaultLogLevel = "warn"
)

func NewEncoderConfig(encoding Encoding) zapcore.EncoderConfig {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}

	if encoding != Console {
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.CallerKey = "caller"
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	}

	return encoderCfg
}

func NewLoggerConfig(logLevel string, encoding Encoding) (*zap.Config, error) {
	level, err := NewLevel(logLevel)

	if err != nil {
		return nil, err
	}

	if encoding != Console && encoding != JSON {
		// transparently convert unknown encodings to default value
		encoding = Default
	}

	encoderCfg := NewEncoderConfig(encoding)

	cfg := &zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Encoding:          string(encoding),
		EncoderConfig:     encoderCfg,
		DisableStacktrace: false,
	}

	return cfg, nil
}

func NewLevel(level string) (zapcore.Level, error) {
	var l zapcore.Level

	err := l.Set(level)

	return l, err
}
