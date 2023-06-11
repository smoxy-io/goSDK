package Logger

import (
	"errors"
	"fmt"
	"github.com/smoxy-io/goSDK/util/logs"
	"go.uber.org/zap"
)

var logger *zap.Logger
var loggers map[string]logs.LoggerInit

func init() {
	loggers = make(map[string]logs.LoggerInit)

	loggers["io"] = InitIOLogger
	loggers["console"] = InitIOLogger
}

func GetLogger() *zap.Logger {
	return logger
}

func RegisterLoggerInit(name string, init logs.LoggerInit) {
	loggers[name] = init
}

func InitLogger(name string, logLevel string, encoding logs.Encoding, options ...zap.Option) (*zap.Logger, error) {
	if logger != nil {
		// logger has already been initialized
		return logger, nil
	}

	init, ok := loggers[name]

	if !ok {
		return nil, errors.New(fmt.Sprintf("no initialization function registered for '%v' logger", name))
	}

	log, err := init(logLevel, encoding, options...)

	if err != nil {
		return nil, err
	}

	logger = log

	return logger, nil
}
