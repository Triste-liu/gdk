package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
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
		return "DEBUG" + strings.Repeat(" ", 2)
	case INFO:
		return "INFO" + strings.Repeat(" ", 3)
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR" + strings.Repeat(" ", 2)
	case PANIC:
		return "PANIC" + strings.Repeat(" ", 2)
	default:
		return "DEBUG" + strings.Repeat(" ", 2)
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
	TraceId  *string
	ClientIp *string
	extra    map[string]interface{}
	skip     bool
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
	st := "\nTraceBack:"
	for {
		frame, more := frames.Next()
		st += fmt.Sprintf("\n%s:%d", strings.Replace(frame.Function, "/", ".", -1), frame.Line)
		if !more {
			break
		}
	}
	return st
}

func getLocation(skip int) string {
	pc, _, line, _ := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc)
	funcName := strings.Replace(fn.Name(), "/", ".", -1)
	return fmt.Sprintf("%s:%d", funcName, line)
}

func writeTextLog(writer io.Writer, r *Record) {
	b := &bytes.Buffer{}
	b.WriteString(r.Level.Color())
	b.WriteString(r.Time.Format(time.DateTime+".000") + " | ")
	b.WriteString(r.Level.String() + " | ")
	b.WriteString(r.Location + " | ")
	if len(r.extra) != 0 {
		e, err := json.Marshal(r.extra)
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
	if r.skip {
		skip--
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
	message = fmt.Sprint(message) + getCaller(3)
	r.log(ERROR, message, args...)
}

func (r *Record) Panic(message interface{}, args ...interface{}) {
	message = fmt.Sprint(message) + getCaller(3)
	r.log(PANIC, message, args...)
	os.Exit(0)
}

func (r *Record) WithFields(e map[string]interface{}) *Record {
	r.extra = e
	r.skip = true
	return r
}
func (r *Record) WithContext(ctx context.Context) *Record {
	traceId, ok := ctx.Value("traceId").(string)
	if ok {
		r.TraceId = &traceId
	}
	clientIp, ok := ctx.Value("clientIp").(string)
	if ok {
		r.ClientIp = &clientIp
	}
	r.skip = true
	return r
}
