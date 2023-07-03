package formatter

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/triste-liu/gdk/log/utils"
)

type ReportFormatter struct {
	TimestampFormat string
}

type JSONLogStruct struct {
	Time     interface{}            `json:"time"`
	Level    logrus.Level           `json:"level"`
	Location string                 `json:"location"`
	Message  string                 `json:"message"`
	Extra    map[string]interface{} `json:"extra"`
}

func (f ReportFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var log JSONLogStruct
	if f.TimestampFormat == "" {
		log.Time = entry.Time.UnixMilli()
	} else {
		log.Time = entry.Time.Format(f.TimestampFormat)
	}
	log.Level = entry.Level
	log.Location = utils.EntryCallerHandler(entry)
	log.Message = entry.Message
	log.Extra = entry.Data

	marshal, err := json.Marshal(log)
	if err != nil {
		return []byte(""), err
	}
	return marshal, nil
}
