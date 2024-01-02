package gjson

import (
	"encoding/json"
	"github.com/triste-liu/gdk/log"
)

func ToString(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		log.Warning(err)
		return ""
	}
	return string(marshal)
}
