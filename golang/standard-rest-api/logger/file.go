package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"runtime"
	"path"
	"sync"
	"archive/zip"
	"io"
	"path/filepath"
	"github.com/golang/standard-rest-api/utils/archive"
)

// fileLogWriter implements LoggerInterface.
// It writes message by file size limit.
type FileLogWriter struct {
	FileName	string
	FilePath 	string
	Level 		int
	LogCount	int
	LogMaxSize	int
	fileWriter	*os.File
	sync.RWMutex
}

// newFileWriter create a FileLogWriter returning as a LoggerInterface
func newFileWriter() *FileLogWriter {
	f := &FileLogWriter{
		FileName:	LogName,
		LogMaxSize:	LogMaxSize,
		LogCount: 	LogCount,
		Level:      	LevelTrace,
	}

	return f
}

// Init init fileLogWriter with config
func (f *FileLogWriter) Init(fileName, filePath string, logLevel, logCount, logMaxSize int) error {
	f.FileName = fileName
	f.FilePath = filePath

	f.setConf(logLevel, logCount, logMaxSize)

	err := f.startLogger()
	if err != nil {
		return err
	}

	go f.updateConfig()
	return nil
}

func (f *FileLogWriter) setConf(logLevel, logCount, logMaxSize int) {
	f.setLogLevel(logLevel)
	f.setLogCount(logCount)
	f.setLogMaxSize(logMaxSize)
}

func (f *FileLogWriter) setLogLevel(logLevel int) {
	if logLevel < LevelTrace || logLevel > LevelCrit {
		f.Level = LevelInfo
		return
	}
	f.Level = logLevel
}

func (f *FileLogWriter) setLogCount(logCount int) {
	if logCount < 0 {
		f.LogCount = LogCount
		return
	}
	f.LogCount = logCount
}

func (f *FileLogWriter) setLogMaxSize(logMaxSize int) {
	if logMaxSize < 0 {
		f.LogMaxSize = LogMaxSize
		return
	}
	f.LogMaxSize = logMaxSize
}

//create log file and set locker-inside file writer
func (f *FileLogWriter) startLogger() error {
	file, err := f.createLogFile()
	if err != nil {
		return err
	}

	if f.fileWriter != nil {
		f.fileWriter.Close()
	}
	f.fileWriter = file

	return nil
}

func (f *FileLogWriter) createLogFile() (*os.File, error) {
	perm, err := strconv.ParseInt("0777", 8, 64)
	if err != nil {
		return nil, err
	}
	fName := f.FilePath + "/" + f.FileName
	fd, err := os.OpenFile(fName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		os.Chmod(fName, os.FileMode(perm))
	}
	return fd, err
}

//writeMsg write logger message into file.
func (f *FileLogWriter) WriteMsg(level int, msg string, v ...interface{}) error {
	if level < f.Level {
		return nil
	}

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}

	logPrefix := f.formatLogPrefix(level)
	msg = logPrefix + msg + "\n"

	f.Lock()
	defer f.Unlock()


		if f.needRotate() {
			if err := f.doRotate(); err != nil {
				fmt.Fprintf(os.Stderr, "doRotate failed, error:%s\n", err)
			}
		}


	if f.fileWriter == nil {
		err = f.startLogger()
	}

	_, err := f.fileWriter.Write([]byte(msg))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write log messsage failed, error:%s\n", err)
	}

	f.doRotateTest()

	return err
}

func (f *FileLogWriter) doRotateTest() error {
	backLog := f.FilePath + "/" + f.FileName + ".1.zip"

	err := f.Compress(f.fileWriter, backLog)
	if err != nil {
		return fmt.Errorf("compress %s\n", err)
	}

	return nil
}


//doRotate rotate the current log file
func (f *FileLogWriter) doRotate() error {
	f.deleteOldLog()

	//rename the log zip
	absPath := f.FilePath + "/" + f.FileName
	for i := f.LogCount - 1; i > 0; i -- {
		//LogName: rest.log.2.zip
		oldLogName := absPath + "." + strconv.Itoa(i) + ".zip"
		newLogName := absPath + "." + strconv.Itoa(i+1) + ".zip"

		if _, err := os.Lstat(oldLogName); err != nil {
			continue
		}

		err := os.Rename(oldLogName, newLogName)
		if err != nil {
			fmt.Errorf("rename %s", err)
		}
	}

	backLog := absPath + ".1.zip"
	if f.fileWriter == nil {
		fmt.Fprintf(os.Stderr, "fileWriter is bad descriptor !\n")
	}

	err := archive.ArchiveZip(absPath, backLog)
	if err != nil {
		return fmt.Errorf("compress %s", err)
	}

	perm, _ := strconv.ParseInt("0660", 8, 64)
	err = os.Chmod(backLog, os.FileMode(perm))

	startLoggerErr := f.restartLogger()
	if startLoggerErr != nil {
		return fmt.Errorf("Rotate restartLogger: %s", startLoggerErr)
	}

	if err != nil {
		return fmt.Errorf("Rotate: %s", err)
	}

	return err
}

func (f *FileLogWriter) restartLogger() error {
	f.fileWriter.Close()

	fName := f.FilePath + "/" + f.FileName
	os.Remove(fName)

	return f.startLogger()
}

func (f *FileLogWriter) deleteOldLog() {
	oldLog := fmt.Sprintf("%s.%d.zip", f.FileName, f.LogCount)
	dir := filepath.Dir(f.FilePath + "/" + f.FileName)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr,
					"Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil {
			return
		}

		if info.Name() == oldLog {
			os.Remove(path)
		}

		return
	})
}

func (f *FileLogWriter) Compress(file *os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	err := f.compress(file, "", w)
	return err
}

func (f *FileLogWriter) compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "compress 1 error:%s", err)
		return err
	}

	header, err := zip.FileInfoHeader(info)
	header.Name = prefix + "/" + header.Name
	if err != nil {
		fmt.Fprintf(os.Stderr, "compress 2 error:%s", err)
		return  err
	}

	writer, err := zw.CreateHeader(header)
	if err != nil {
		fmt.Fprintf(os.Stderr, "compress 3 error:%s", err)
		return err
	}

	if file == nil {
		fmt.Fprintf(os.Stderr, "file *os.File is nill")
	}

	_, err = io.Copy(writer, file)
	fmt.Fprintf(os.Stderr, "compress 4 error:%s", err)
	return err
}

func (f *FileLogWriter) needRotate() bool {
	fileInfo, err := f.fileWriter.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "get state error:%s\n", err)
		return false
	}

	if f.LogMaxSize < int(fileInfo.Size()) {
		fmt.Fprintf(os.Stderr, "logMaxSize:%d, fileSize:%d\n", f.LogMaxSize, int(fileInfo.Size()))
		return true
	}

	return false
}

func (f *FileLogWriter) formatLogPrefix(level int) string {
	when := time.Now()
	h, _ := formatTimeHeader(when)

	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)

	return fmt.Sprintf("[%s][%s][%s,%s]",
		string(h), levelPrefix[level], filename, strconv.Itoa(line))
}


func (f *FileLogWriter) updateConfig() {
	//TODO: update info from the log config
	for {
		time.Sleep(1 * time.Second)
	}
}