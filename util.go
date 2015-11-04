package gochat

import (
	"io"
	"log"
	"os"
)

var (
	LTrace   *log.Logger
	LWarning *log.Logger
	LError   *log.Logger
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func LogInit(
	traceHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	LTrace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LWarning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LError = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
