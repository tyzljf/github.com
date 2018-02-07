package logs

import (
	"os"
	"time"
	"sync"
	"fmt"
)

//log message levels.
const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	CRIT
)

const (
	DEFAULT_LOG_NAME = "standard-rest.log"
	DEFAULT_LOG_SIZE = 30 * 1024 * 1024 //30M
	DEFAULT_LOG_COUNT = 3
)

//implement logger instance singleton
var tlLogger *Logger

func NewLogger() *Logger {
	tlLogger = &Logger{
		Level: INFO,
		FileName: DEFAULT_LOG_NAME,
		LogCount: 3,
		LogMaxSize: DEFAULT_LOG_SIZE,
	}
	return tlLogger
}

type Logger struct {
	Level int
	FileName string
	FilePath string
	LogCount int
	LogMaxSize int
	lock sync.Mutex
	fileWriter *os.File
}

func (tl *Logger) Init(fileName, filePath string, logLevel, logCount, logMaxSize int) {
	tl.FileName = fileName
	tl.FilePath = filePath
	tl.setConf(logLevel, logCount, logMaxSize)
}

func (tl *Logger) setConf(logLevel, logCount, logMaxSize int) {
	tl.setLogLevel(logLevel)
	tl.setLogCount(logCount)
	tl.setLogMaxSize(logMaxSize)
}

func (tl *Logger) setLogLevel(logLevel int) {
	if logLevel < TRACE || logLevel > CRIT {
		tl.Level = INFO
		return
	}
	tl.Level = logLevel
}

func (tl *Logger) setLogCount(logCount int) {
	if logCount < 0 {
		tl.LogCount = DEFAULT_LOG_SIZE
		return
	}
	tl.LogCount = logCount
}

func (tl *Logger) setLogMaxSize(logMaxSize int) {
	if logMaxSize < 0 {
		tl.LogMaxSize = DEFAULT_LOG_SIZE
		return
	}
	tl.LogMaxSize = logMaxSize
}

func (tl *Logger) writeMsg(level int, msg string, v ...interface{}) error {
	if level > tl.Level {
		return nil
	}

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v)
	}

	//compose the log message
	head := tl.makeHead(level)
}

func (tl *Logger) makeHead(level int) string {
	head := ""

	switch level {
	case TRACE:
		head := "TRACE"
		break
	case INFO:

	}
}

func Trace(msg string, v ...interface{}) {
	tlLogger.writeMsg(TRACE, msg, v)
}

func Debug(msg string, v ...interface{}) {
	tlLogger.writeMsg(DEBUG, msg, v)
}

func Info(msg string, v ...interface{}) {
	tlLogger.writeMsg(INFO, msg, v)
}

func Warn(msg string, v ...interface{}) {
	tlLogger.writeMsg(WARN, msg, v)
}

func Crit(msg string, v ...interface{}) {
	tlLogger.writeMsg(CRIT, msg, v)
}






















