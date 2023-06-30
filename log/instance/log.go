package instance

import (
	"github.com/Triste-liu/gdk/log/formatter"
	"github.com/Triste-liu/gdk/log/hook"
	"github.com/sirupsen/logrus"
	"io"
)

type Config struct {
	ReportWriter    []io.Writer
	ReportLevel     logrus.Level
	ReportFormatter logrus.Formatter
	TraceLevel      logrus.Level
}

func New(config Config) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(formatter.Formatter{})
	log.SetLevel(logrus.TraceLevel)
	log.SetReportCaller(true)
	traceHook := hook.StackHook{Level: config.TraceLevel}
	log.AddHook(&traceHook)
	reportHook := hook.ReportHook{Level: config.ReportLevel, Writer: config.ReportWriter, Formatter: config.ReportFormatter}
	log.AddHook(&reportHook)
	log.Debug("日志初始化成功")
	return log
}

//func main() {
//	httpWriter := &writer.HttpWriter{
//		Url:    "http://localhost:4567/api/v2/log",
//		Method: "POST",
//	}
//	httpHook := hook.ReportHook{Level: logrus.TraceLevel, Writer: httpWriter, Formatter: format.ReportFormatter{}}
//
//	log := New()
//	log.AddHook(&httpHook)
//	log.Info("hello")
//	log.Debug("hello")
//	log.Trace("hello")
//}
