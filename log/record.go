package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	PANIC
)

const (
	blue   = "\033[97;34m"
	green  = "\033[97;32m"
	yellow = "\033[90;33m"
	red    = "\033[97;31m"
	reset  = "\033[0m"
)

type Level int

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG  "
	case INFO:
		return "INFO   "
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR  "
	case PANIC:
		return "PANIC  "
	default:
		return "DEBUG  "
	}
}

func (l Level) Color() string {
	switch l {
	case DEBUG:
		return blue
	case INFO:
		return green
	case WARNING:
		return yellow
	case ERROR:
		return red
	case PANIC:
		return red
	default:
		return blue
	}
}

type Record struct {
	Time     time.Time
	Level    Level
	Location string
	Message  string
	Extra    map[string]interface{}
}

func (r *Record) Byte() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return b
}

func getCaller(skip int) string {
	p := make([]uintptr, 32)
	n := runtime.Callers(skip, p)
	var frames *runtime.Frames
	if n > 2 {
		frames = runtime.CallersFrames(p[:n-2])
	} else {
		frames = runtime.CallersFrames(p)
	}
	var st string
	for {
		frame, more := frames.Next()
		st += fmt.Sprintf("\n%s:%d  ->  %s", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	return st
}

func getLocation(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc)
	///go/src/control-center/service/chat/handler.go:chat.abnormal:77
	return fmt.Sprintf("%s:%s:%d", file, fn.Name(), line)
}

func writeTextLog(writer io.Writer, r *Record) {
	b := &bytes.Buffer{}
	b.WriteString(r.Level.Color())
	b.WriteString(r.Time.Format(time.DateTime+".000") + " | ")
	b.WriteString(r.Level.String() + " | ")
	b.WriteString(r.Location + " | ")
	if len(r.Extra) != 0 {
		e, err := json.Marshal(r.Extra)
		if err != nil {
			fmt.Printf("extra marshal error:%s\n", err)
		} else {
			b.Write(e)
			b.WriteString(" | ")
		}

	}
	b.WriteString(r.Message)
	b.WriteString(reset + "\n")
	_, err := writer.Write(b.Bytes())
	if err != nil {
		fmt.Printf("write log error: %v\n", err)
	}
}

func writeJsonLog(writer io.Writer, r *Record) {
	_, err := writer.Write(r.Byte())
	if err != nil {
		fmt.Printf("write log error: %v\n", err)
	}
}

func (r *Record) log(level Level, message interface{}, args ...interface{}) {
	r.Level = level
	r.Time = time.Now()
	if len(args) == 0 {
		r.Message = fmt.Sprint(message)
	} else {
		r.Message = fmt.Sprintf(fmt.Sprint(message), args...)
	}
	skip := 4
	if len(r.Extra) != 0 {
		skip++
	}
	r.Location = getLocation(skip)
	for _, v := range instance {
		if r.Level >= v.Level {
			switch v.Type {
			case TEXT:
				writeTextLog(v.Writer, r)
			case JSON:
				writeJsonLog(v.Writer, r)
			default:
				writeTextLog(v.Writer, r)
			}
		}
	}
}

func (r *Record) Debug(message interface{}, args ...interface{}) {
	r.log(DEBUG, message, args...)
}

func (r *Record) Info(message interface{}, args ...interface{}) {
	r.log(INFO, message, args...)
}

func (r *Record) Warning(message interface{}, args ...interface{}) {
	r.log(WARNING, message, args...)
}

func (r *Record) Error(message interface{}, args ...interface{}) {
	message = fmt.Sprint(message) + getCaller(4)
	r.log(ERROR, message, args...)
}

func (r *Record) Panic(message interface{}, args ...interface{}) {
	message = fmt.Sprint(message) + getCaller(4)
	r.log(PANIC, message, args...)
	os.Exit(0)
}
