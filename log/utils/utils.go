package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func EntryCallerHandler(entry *logrus.Entry) (caller string) {
	funcSplit := strings.Split(entry.Caller.Function, "/")
	function := funcSplit[len(funcSplit)-1]
	caller = fmt.Sprintf("%s:%s:%d", entry.Caller.File, function, entry.Caller.Line)
	return
}
