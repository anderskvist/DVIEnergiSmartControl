package log

import (
	"os"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

// Debug blah blah blah
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Fatal blah blah blah
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Info blah blah blah
func Info(args ...interface{}) {
	log.Info(args...)
}

// Notice blah blah blah
func Notice(args ...interface{}) {
	log.Notice(args...)
}

// Infof blah blah blah
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Noticef blah blah blah
func Noticef(format string, args ...interface{}) {
	log.Noticef(format, args...)
}
