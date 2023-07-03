package formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/triste-liu/gdk/log/utils"
	"sort"
	"strings"
	"time"
)

const (
	green   = "\033[97;32m"
	white   = "\033[90;37m"
	yellow  = "\033[90;33m"
	red     = "\033[97;31m"
	blue    = "\033[97;34m"
	magenta = "\033[97;35m"
	cyan    = "\033[97;36m"
	reset   = "\033[0m"
)

type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.DateTime   = "2006-01-02 15:04:05"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsSpace - no space between fields
	FieldsSeparator string
}

func getColorByLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return blue
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	default:
		return green
	}
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		// a.b:c
		// E:/service-center/service/service.go  service-center/service.PingJob   59
		//fmt.Fprintf(
		//	b,
		//	"%s:%d",
		//	entry.Caller.Function,
		//	entry.Caller.Line,
		//)
		caller := utils.EntryCallerHandler(entry)
		b.WriteString(caller)
	}
	b.WriteString(f.FieldsSeparator)
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}
	b.WriteString(f.FieldsSeparator)
}
func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.FieldsSeparator == "" {
		f.FieldsSeparator = " | "
	}
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.DateTime + ".000" //milli
	}

	// output buffer
	b := &bytes.Buffer{}

	if !f.NoColors {
		b.WriteString(levelColor)
	}
	// write time
	b.WriteString(entry.Time.Format(timestampFormat) + f.FieldsSeparator)

	// write level
	var level string
	level = strings.ToUpper(entry.Level.String())
	b.WriteString("[" + level + "]" + f.FieldsSeparator)

	// write caller
	if entry.HasCaller() {
		f.writeCaller(b, entry)
	}

	// write fields
	f.writeFields(b, entry)

	// write message
	b.WriteString(entry.Message)

	b.WriteByte('\n')

	// reset color
	if !f.NoColors {
		b.WriteString(reset)
	}

	return b.Bytes(), nil
}
