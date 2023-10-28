package log

import (
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

func Extra(e map[string]interface{}) *Record {
	r := &Record{Extra: e}
	return r
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
