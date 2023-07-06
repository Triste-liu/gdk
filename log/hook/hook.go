package hook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

type ReportHook struct {
	Level     logrus.Level
	Writer    []io.Writer
	Formatter logrus.Formatter
}

func (h *ReportHook) Levels() []logrus.Level {
	return logrus.AllLevels[:h.Level+1]
}

func (h *ReportHook) Fire(entry *logrus.Entry) error {
	bytes, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}
	for _, writer := range h.Writer {
		_, err := writer.Write(bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

type StackHook struct {
	Level logrus.Level
}

func (h *StackHook) Levels() []logrus.Level {
	return logrus.AllLevels[:h.Level+1]
}

func (h *StackHook) Fire(entry *logrus.Entry) error {
	p := make([]uintptr, 32)
	n := runtime.Callers(8, p)
	frames := runtime.CallersFrames(p[:n])
	for {
		frame, more := frames.Next()
		st := fmt.Sprintf("\n%s:%d  ->  %s", frame.File, frame.Line, frame.Function)
		entry.Message += st
		if !more {
			break
		}
	}
	return nil
}
