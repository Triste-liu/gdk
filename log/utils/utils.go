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
	goPath := os.Getenv("GOPATH")
	goPath = strings.Replace(goPath, "\\", "/", -1)
	goPath = goPath + "/pkg/mod/"

	parentDir := filepath.Dir(filepath.Clean(wd))
	parentDir = strings.Replace(parentDir, "\\", "/", -1)
	file := strings.TrimPrefix(entry.Caller.File, parentDir+"/")
	file = strings.TrimPrefix(file, goPath)
	funcSplit := strings.Split(entry.Caller.Function, "/")
	function := funcSplit[len(funcSplit)-1]
	caller = fmt.Sprintf("%s:%s:%d", file, function, entry.Caller.Line)
	return
}
