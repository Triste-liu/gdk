package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func EntryCallerHandler(entry *logrus.Entry) (caller string) {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	parentDir := filepath.Dir(filepath.Clean(wd))
	parentDir = strings.Replace(parentDir, "\\", "/", -1)
	file := strings.TrimPrefix(entry.Caller.File, parentDir+"/")
	funcSplit := strings.Split(entry.Caller.Function, "/")
	function := funcSplit[len(funcSplit)-1]
	caller = fmt.Sprintf("%s:%s:%d", file, function, entry.Caller.Line)
	return
}
