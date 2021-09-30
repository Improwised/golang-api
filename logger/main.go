package logger

import (
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var zapDebugConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Development:      true,
	Encoding:         "console",
	EncoderConfig:    defaultEncoderConfig,
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

var zapServerConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	Development:      false,
	Encoding:         "json",
	EncoderConfig:    defaultEncoderConfig,
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

//NewRootLogger instantiates zap.Logger with given configuration
func NewRootLogger(debug, developement bool) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger
	if debug {
		if !developement {
			zapDebugConfig.Encoding = "json"
			return zapDebugConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
		}
		zapDebugConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapDebugConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	if developement {
		zapDebugConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapDebugConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	log.SetFormatter(&log.JSONFormatter{})
	logger, err = zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())

	if err != nil {
		panic(err)
	}

	return logger, err
}
