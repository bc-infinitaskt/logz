package logz

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = zap.NewExample()

type Option func()

func WithLogEnv(env string) Option {
	return func() {
		var cfg zap.Config
		cfg = zap.NewProductionConfig()
		switch env {
		case "prod":
			cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		case "debug":
			cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		default:
			cfg = zap.NewDevelopmentConfig()
		}

		cfg.EncoderConfig = zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "@timestamp",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		}

		encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
		core := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), cfg.Level)
		*logger = *zap.New(core, zap.AddCaller())
	}
}

func WithAppName(appName string) Option {
	return func() {
		*logger = *logger.With(zap.String("app", appName))
	}
}

func NewLogger(opts ...Option) *zap.Logger {
	for _, opt := range opts {
		opt()
	}
	return logger
}
