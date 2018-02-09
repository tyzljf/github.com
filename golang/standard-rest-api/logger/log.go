package logger

import (
	"time"
	"fmt"
	"os"
	"errors"
)

//log message levels
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelCrit
)

const (
	LogName = "standard-rest.log"
	LogMaxSize = 30 * 1024 * 1024 //30M
	LogCount = 3
)

type Logger interface {
	Init(fileName, filePath string, logLevel, logCount, logMaxSize int) error
	WriteMsg(when time.Time, msg string, level int) error
}

func Init(fileName, filePath string, logLevel, logCount, logMaxSize int) error {
	if len(fileName) == 0 {
		return errors.New("filename is empty")
	}

	err := tlLogger.Init(fileName, filePath, logLevel, logCount, logMaxSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logger failed, error:%s", err)
	}

	return err
}

var tlLogger = newFileWriter()
var levelPrefix = [LevelCrit + 1]string{"TRACE", "DEBUG", "INFO", "WARN", "CRIT"}

func Trace(format string, v ...interface{}) {
	tlLogger.WriteMsg(LevelTrace, format, v...)
}

func Debug(format string, v ...interface{}) {
	tlLogger.WriteMsg(LevelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
	tlLogger.WriteMsg(LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	tlLogger.WriteMsg(LevelWarn, format, v...)
}

func Crit(format string, v ...interface{}) {
	tlLogger.WriteMsg(LevelCrit, format, v...)
}

const (
	y1  = `0123456789`
	y2  = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789`
	y3  = `0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999`
	y4  = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789`
	mo1 = `000000000111`
	mo2 = `123456789012`
	d1  = `0000000001111111111222222222233`
	d2  = `1234567890123456789012345678901`
	h1  = `000000000011111111112222`
	h2  = `012345678901234567890123`
	mi1 = `000000000011111111112222222222333333333344444444445555555555`
	mi2 = `012345678901234567890123456789012345678901234567890123456789`
	s1  = `000000000011111111112222222222333333333344444444445555555555`
	s2  = `012345678901234567890123456789012345678901234567890123456789`
	ns1 = `0123456789`
)

func formatTimeHeader(when time.Time) ([]byte, int) {
	y, mo, d := when.Date()
	h, mi, s := when.Clock()
	ns := when.Nanosecond()/1000000
	//len("2006/01/02 15:04:05.123 ")==24
	var buf [24]byte

	buf[0] = y1[y/1000%10]
	buf[1] = y2[y/100]
	buf[2] = y3[y-y/100*100]
	buf[3] = y4[y-y/100*100]
	buf[4] = '/'
	buf[5] = mo1[mo-1]
	buf[6] = mo2[mo-1]
	buf[7] = '/'
	buf[8] = d1[d-1]
	buf[9] = d2[d-1]
	buf[10] = ' '
	buf[11] = h1[h]
	buf[12] = h2[h]
	buf[13] = ':'
	buf[14] = mi1[mi]
	buf[15] = mi2[mi]
	buf[16] = ':'
	buf[17] = s1[s]
	buf[18] = s2[s]
	buf[19] = '.'
	buf[20] = ns1[ns/100]
	buf[21] = ns1[ns%100/10]
	buf[22] = ns1[ns%10]

	buf[23] = ' '

	return buf[0:23], d
}