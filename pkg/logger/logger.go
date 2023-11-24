package logger

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var DefaultZapOptions []zap.Option = []zap.Option{
	zap.AddCaller(),
	// zap.AddStacktrace(zapcore.ErrorLevel),
	zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSamplerWithOptions(
			core,
			time.Second,
			100,
			5,
		)
	}),
}

var JsonEncoder zapcore.Encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
})

var ConsoleEncoder zapcore.Encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
})

func NewLogger(config Config, encoder zapcore.Encoder, options ...zap.Option) (*zap.Logger, error) {
	switch config.Type {
	case File:
		return NewFileLogger(config, encoder, options...)
	case Console:
		return NewConsoleLogger(config, encoder, options...), nil
	default:
		return NewConsoleLogger(config, encoder, options...), nil
	}
}

func NewFileLogger(config Config, encoder zapcore.Encoder, options ...zap.Option) (*zap.Logger, error) {
	roller, err := lumberjack.NewRoller(
		config.LogPath,
		config.MaxSize,
		&lumberjack.Options{
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		},
	)
	if err != nil {
		return nil, err
	}

	logger := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(roller), GetLevel(config.Level)),
		options...,
	)

	return logger, nil
}

func NewConsoleLogger(config Config, encoder zapcore.Encoder, options ...zap.Option) *zap.Logger {
	return zap.New(
		zapcore.NewCore(encoder, os.Stdout, GetLevel(config.Level)),
		options...,
	)
}
