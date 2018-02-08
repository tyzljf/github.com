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
var tlSync   sync.Mutex

func init() {
	tlSync.Lock()
	defer tlSync.Unlock()
	if tlLogger == nil {
		tlLogger = &Logger{
			Level: INFO,
			FileName: DEFAULT_LOG_NAME,
			LogCount: DEFAULT_LOG_COUNT,
			LogMaxSize: DEFAULT_LOG_SIZE,
		}
	}
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
	if level < tl.Level {
		return nil
	}

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}

	//compose the log message
	head := tl.formatHead(level)
	time_str := tl.formatTime()

	msg = fmt.Sprintf("[%s][%s]%s\n", time_str, head, msg)
	return tl.writeToLogger(msg)
}

func (tl *Logger) writeToLogger(msg string) error {
	if tl.fileWriter != nil {
		tl.fileWriter.Close()
	}

	file, err := os.OpenFile(tl.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	tl.fileWriter = file

	_, err =  tl.fileWriter.Write([]byte(msg))
	return err
}

func (tl *Logger) formatTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func (tl *Logger) formatHead(level int) string {
	head := " "

	switch level {
	case TRACE:
		head = "TRACE"
		break
	case DEBUG:
		head = "DEBUG"
		break
	case INFO:
		head = "INFO"
		break
	case WARN:
		head = "WARN"
		break
	case CRIT:
		head = "CRIT"
		break
	}

	return head
}

func Trace(format string, v ...interface{}) {
	tlLogger.writeMsg(TRACE, format, v...)
}

func Debug(format string, v ...interface{}) {
	tlLogger.writeMsg(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	tlLogger.writeMsg(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	tlLogger.writeMsg(WARN, format, v...)
}

func Crit(format string, v ...interface{}) {
	tlLogger.writeMsg(CRIT, format, v...)
}






















