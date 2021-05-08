package log

import (
	"context"
	"fmt"
)

var X *Logger

type method int8

const (
	Debug method = iota - 1
	Info
	Warn
	Error
)

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
		value := ctx.Value("_log")
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
