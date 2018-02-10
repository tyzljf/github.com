package env

import (
	"path/filepath"
	"os"
	"fmt"
	"strings"
)

const (
	PathSeparator = "/"
)

var RootDir = GetCurrentDir()

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "get current directory failed, error:%s", err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

func GetLogPath() string {
	return RootDir + PathSeparator + "logs"
}