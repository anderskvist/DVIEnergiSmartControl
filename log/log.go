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
	cfg, err := ini.Load(os.Args[1])

	if err != nil {
		log.Criticalf("Fail to read file: %v", err)
		os.Exit(1)
	}

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

// Critical blah blah blah
func Critical(args ...interface{}) {
	log.Critical(args...)
}

// Debug blah blah blah
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Error blah blah blah
func Error(args ...interface{}) {
	log.Error(args...)
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

// Criticalf blah blah blah
func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}

// Debugf blah blah blah
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Errorf blah blah blah
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf blah blah blah
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Infof blah blah blah
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Noticef blah blah blah
func Noticef(format string, args ...interface{}) {
	log.Noticef(format, args...)
}
