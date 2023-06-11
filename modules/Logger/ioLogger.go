package Logger

import (
	"github.com/smoxy-io/goSDK/util/logs"
	"go.uber.org/zap"
)

func InitIOLogger(logLevel string, encoding logs.Encoding, options ...zap.Option) (*zap.Logger, error) {
	cfg, err := logs.NewLoggerConfig(logLevel, encoding)

	if err != nil {
		return nil, err
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	log, err := cfg.Build(options...)

	return log, err
}
