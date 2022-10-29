package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/google/wire"
	"github.com/natefinch/lumberjack/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	logger *zap.Logger
	Config logger.Config
}

func GetLogLevel(level config.LogLevel) zapcore.Level {
	var l zapcore.Level

	switch level {
	case config.Debug:
		l = zapcore.DebugLevel
	case config.Info:
		l = zapcore.InfoLevel
	case config.Warn:
		l = zapcore.WarnLevel
	case config.Error:
		l = zapcore.ErrorLevel
	case config.Dpanic:
		l = zapcore.DPanicLevel
	case config.Panic:
		l = zapcore.PanicLevel
	case config.Fatal:
		l = zapcore.FatalLevel
	case config.Silent:
		l = zapcore.FatalLevel
	}
	return l
}

func GetGormLogLevel(level config.LogLevel) logger.LogLevel {
	var l logger.LogLevel

	switch level {
	case config.Debug:
		l = logger.Info
	case config.Info:
		l = logger.Info
	case config.Warn:
		l = logger.Warn
	case config.Error:
		l = logger.Error
	case config.Dpanic:
		l = logger.Silent
	case config.Panic:
		l = logger.Silent
	case config.Fatal:
		l = logger.Silent
	case config.Silent:
		l = logger.Silent
	}
	return l
}

var zapOptions []zap.Option = []zap.Option{
	zap.AddCaller(),
	// zap.AddStacktrace(zapcore.ErrorLevel),
	zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSamplerWithOptions(
			core,
			time.Second,
			100,
			100,
		)
	}),
}

var fileEncoder zapcore.Encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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

var logConfig = config.Log()

func NewLogger(writer io.Writer) *zap.Logger {
	var writeSyncer zapcore.WriteSyncer

	switch logConfig.Type {
	case config.Console:
		writeSyncer = os.Stdout
	case config.File:
		writeSyncer = zapcore.AddSync(writer)
	}
	zapLogger := zap.New(
		zapcore.NewCore(fileEncoder, writeSyncer, GetLogLevel(logConfig.Level)),
		zapOptions...,
	)

	return zapLogger
}

func NewAccessLogger() *zap.Logger {
	roller, err := lumberjack.NewRoller(
		path.Join(config.App().ProjectRoot, logConfig.AccessLogPath),
		logConfig.MaxSize,
		&lumberjack.Options{
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		},
	)
	if err != nil {
		panic(err)
	}

	return NewLogger(roller)
}

func provideLogger() *zap.Logger {
	roller, err := lumberjack.NewRoller(
		path.Join(config.App().ProjectRoot, logConfig.AppLogPath),
		logConfig.MaxSize,
		&lumberjack.Options{
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		},
	)
	if err != nil {
		panic(err)
	}

	return NewLogger(roller)
}

func provideGormLogger(zapLogger *zap.Logger) *gormLogger {
	return &gormLogger{
		logger: zapLogger,
		Config: logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  GetGormLogLevel(logConfig.Level),
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
	}
}

var logProviderSet = wire.NewSet(provideLogger, provideGormLogger)

func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *gl
	newlogger.Config.LogLevel = level

	return &newlogger
}

func (gl *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if gl.Config.LogLevel >= logger.Info {
		gl.logger.Info(fmt.Sprintf(msg, data...))
	}
}

func (gl *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if gl.Config.LogLevel >= logger.Warn {
		gl.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (gl *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if gl.Config.LogLevel >= logger.Error {
		gl.logger.Error(fmt.Sprintf(msg, data...))
	}
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.Config.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && (!errors.Is(err, logger.ErrRecordNotFound) || !gl.Config.IgnoreRecordNotFoundError):
		if rows == -1 {
			gl.logger.Error(fmt.Sprintf("%s %s %fms %s %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			gl.logger.Error(fmt.Sprintf("%s %s %fms [rows:%d] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	case elapsed > gl.Config.SlowThreshold && gl.Config.SlowThreshold != 0:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", gl.Config.SlowThreshold)
		if rows == -1 {
			gl.logger.Warn(fmt.Sprintf("%s %s %fms %s %s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			gl.logger.Warn(fmt.Sprintf("%s %s %fms [rows:%d] %s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	default:
		if rows == -1 {
			gl.logger.Debug(fmt.Sprintf("%s %fms %s %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			gl.logger.Debug(fmt.Sprintf("%s %fms [rows:%d] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	}
}
