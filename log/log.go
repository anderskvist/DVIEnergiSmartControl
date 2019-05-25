package log

import (
	"os"

	logging "github.com/op/go-logging"
	ini "gopkg.in/ini.v1"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	cfg, _ := ini.Load(os.Args[1])
	loglevel := cfg.Section("main").Key("loglevel").String()
	level, _ := logging.LogLevel(loglevel)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)

	Noticef("Enabling logging at loglevel: %s\n", logging.GetLevel("example"))
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
