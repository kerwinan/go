package log

import (
	"context"
	"fmt"
)

var X *Logger

const contextLogKey = "_log"

type method int8

const (
	Debug method = iota - 1
	Info
	Warn
	Error
)

// Flush calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (it *Logger) Flush() {
	it.zapLogger.Sync()
}

// With adds entries and constructs a new Logger.
// Note that the keys in key-value pairs should be strings.
func (it *Logger) With(keysAndValues ...interface{}) (log *Logger) {
	log = new(Logger)
	log.level = it.level
	log.sugar = it.sugar.With(keysAndValues...)
	return
}

// Withc adds entries and constructs a new Logger, and uses fmt.Sprintf to store a templated message.
func (it *Logger) Withf(key string, format string, params ...interface{}) (log *Logger) {
	log = new(Logger)
	log.level = it.level
	log.sugar = it.sugar.With(key, fmt.Sprintf(format, params...))
	return
}

//  same as With, but store in context
func (it *Logger) Withc(ctx context.Context, keysAndValues ...interface{}) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	value := ctx.Value(contextLogKey)
	var kvs []interface{}
	switch value.(type) {
	case []interface{}:
		kvs = append(value.([]interface{}), keysAndValues...)
	default:
		kvs = keysAndValues
	}
	return context.WithValue(ctx, contextLogKey, kvs)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (it *Logger) Debugf(format string, params ...interface{}) {
	it.generate(nil, it, Debug, format, params...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (it *Logger) Infof(format string, params ...interface{}) {
	it.generate(nil, it, Info, format, params...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (it *Logger) Warnf(format string, params ...interface{}) {
	it.generate(nil, it, Warn, format, params...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (it *Logger) Errorf(format string, params ...interface{}) {
	it.generate(nil, it, Error, format, params...)
}

func (it *Logger) generate(ctx context.Context, self *Logger, fun method, format interface{}, params ...interface{}) {
	if it.level > fun {
		return
	}

	var msg string
	if len(format.(string)) > 0 {
		msg = fmt.Sprintf(format.(string), params...)
	} else {
		for i, param := range params {
			if i == 0 {
				msg += fmt.Sprintf("%+v", param)
			} else {
				msg += fmt.Sprintf(" %+v", param)
			}
		}
	}
	var sugar = self.sugar

	if ctx != nil {
		value := ctx.Value(contextLogKey)
		switch value.(type) {
		case []interface{}:
			sugar = sugar.With(value.([]interface{})...)
		default:
		}
	}

	switch fun {
	case Debug:
		sugar.Debug(msg)
	case Info:
		sugar.Info(msg)
	case Warn:
		sugar.Warn(msg)
	case Error:
		sugar.Error(msg)
	}
	return
}
