package log

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"io"
)

// Config is struct about configuration file.
type Config struct {
	// Level is the minimum log level
	// -1 is debug level,
	// 0 is info level,
	// 1 is warn level,
	// 2 is error level
	Level int8 `yaml:"level"`

	// keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey string `yaml:"message_key"` // key of message
	LevelKey   string `yaml:"level_key"`   // key of level
	TimeKey    string `yaml:"time_key"`    // key of time
	CallerKey  string `yaml:"caller_key"`  // key of caller

	// encoding of log, just is json or console
	Encoding string `yaml:"encoding"`

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string `yaml:"file_name"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `yaml:"local_time"`

	writer io.Writer
}

//Option option
type Option func(conf *Config)

type Log interface {
	Flush()
	With(...interface{}) *Logger
	Withf(string, string, ...interface{}) *Logger
	Withc(context.Context, ...interface{}) context.Context
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

// Logger is the implement about Log.
type Logger struct {
	zapLogger  *zap.Logger
	sugar *zap.SugaredLogger
	lumberjack *lumberjack.Logger
	level method
}
