package logs

import (
	"encoding/json"
	"errors"
	"os"
)

//fileLogWriter implements LoggerInterface
type fileLogWriter struct {
	Filename string	`json:"filename"`
	fileWriter *os.File
}

func newFileWriter() *Logger {
	w := &fileLogWriter{

	}
	return w
}

//Init file logger with json config
func (w *fileLogWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), w)
	if err != nil {
		return err
	}
	if len(w.Filename) == 0 {
		return errors.New("jsonconfig must have filename")
	}
	err = w.startLogger()
	return err
}

func (w *fileLogWriter) startLogger() error {
	file, err := w.createLogFile()
	if err != nil {
		return err
	}
	if w.fileWriter != nil {
		w.fileWriter.Close()
	}
	w.fileWriter = file
	return w.initFd()
}