package logger

import (
	"github.com/IceFoxs/open-gateway/conf"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *hertzzap.Logger {
	log := conf.GetConf().Logger
	// 提供压缩和删除
	zapLoggerEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         log.Encoding,
		EncoderConfig:    zapLoggerEncoderConfig,
		OutputPaths:      []string{"stdout", log.FileName},
		ErrorOutputPaths: []string{"stderr", log.FileName},
	}
	lumberjackLogger := &lumberjack.Logger{
		Filename:   log.FileName,
		MaxSize:    log.MaxSize,
		MaxBackups: log.MaxBackups,
		MaxAge:     log.MaxAge,
		Compress:   true,
	}
	logger := hertzzap.NewLogger(hertzzap.WithZapOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			getEncoder(config),
			getLogWriter(lumberjackLogger),
			config.Level)
	})))
	return logger
}

func getEncoder(z *zap.Config) zapcore.Encoder {
	if z.Encoding == "json" {
		return zapcore.NewJSONEncoder(z.EncoderConfig)
	} else if z.Encoding == "console" {
		return zapcore.NewConsoleEncoder(z.EncoderConfig)
	}
	return nil
}

// getLogWriter get Lumberjack writer by LumberjackConfig
func getLogWriter(l *lumberjack.Logger) zapcore.WriteSyncer {
	return zapcore.AddSync(l)
}
