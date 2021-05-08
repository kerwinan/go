package log

import (
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// GetLoggerByConf constructs a new Logger by Config.
func GetLoggerByConf(config *Config, options ...Option) (logger *Logger, err error) {
	for _, opt := range options {
		opt(config)
	}

	proConf := zapcore.EncoderConfig{
		MessageKey:     config.MessageKey,
		LevelKey:       config.LevelKey,
		TimeKey:        config.TimeKey,
		CallerKey:      config.CallerKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// choose the type of encoding.
	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(proConf)
	} else if config.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(proConf)
	} else {
		err = errors.New("encoding must be one of the json, console or simple")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var lumberjackLogger *lumberjack.Logger
	var zapWriter, output zapcore.WriteSyncer
	if config.writer != nil {
		// write logs to custom writer
		zapWriter = zapcore.AddSync(config.writer)
		output = zapcore.NewMultiWriteSyncer(zapWriter)
	} else {
		// write logs to rolling files
		lumberjackLogger = &lumberjack.Logger{
			Filename:      config.Filename,
			LocalTime:     config.LocalTime,
		}
		zapWriter = zapcore.AddSync(lumberjackLogger)
		output = zapcore.NewMultiWriteSyncer(zapWriter)
		if config.Filename == "" {
			output = os.Stdout
		}
	}

	newCore := zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(zapcore.Level(config.Level)))
	opts := []zap.Option{zap.ErrorOutput(zapWriter)}
	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(2))

	logger = new(Logger)
	logger.zapLogger = zap.New(newCore, opts...)
	logger.lumberjack = lumberjackLogger
	logger.sugar = logger.zapLogger.Sugar()
	logger.level = method(config.Level)
	return
}

// timeEncoder sets the time format
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
