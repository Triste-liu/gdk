package log

import (
	"context"
	"io"
	"os"
)

type Type int

const (
	JSON = iota
	TEXT
)

type Logger struct {
	Type   Type
	Level  Level
	Writer io.Writer
}

var instance = []Logger{{Level: DEBUG, Writer: os.Stdout, Type: TEXT}}

var record Record

func Add(l Logger) {
	instance = append(instance, l)
}

func SetLevel(level Level) {
	instance[0].Level = level
}

func WithFields(e map[string]interface{}) *Record {
	return &Record{extra: e, skip: true}
}

func Debug(message interface{}, args ...interface{}) {
	record.Debug(message, args...)
}

func Info(message interface{}, args ...interface{}) {
	record.Info(message, args...)
}

func Warning(message interface{}, args ...interface{}) {
	record.Warning(message, args...)
}

func Error(message interface{}, args ...interface{}) {
	record.Error(message, args...)
}

func Panic(message interface{}, args ...interface{}) {
	record.Panic(message, args...)
}

func WithContext(ctx context.Context) *Record {
	var r Record
	traceId, ok := ctx.Value("traceId").(string)
	if ok {
		r.TraceId = &traceId
	}
	clientIp, ok := ctx.Value("clientIp").(string)
	if ok {
		r.ClientIp = &clientIp
	}
	r.skip = true
	return &r
}
