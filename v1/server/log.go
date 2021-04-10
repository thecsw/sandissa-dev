package main

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// log is a global logrus instance.
	log = logrus.New()
)

// params is an alias for `map[string]interface{}`.
type params map[string]interface{}

// linit initializes the global logrus instance.
func linit() {
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
		DisableSorting:  false,
	})
	log.SetOutput(os.Stdout)
	l("Logger is created")
}

// lf is an INFO log output with fields.
func lf(what string, fields params) {
	log.WithFields(logrus.Fields(fields)).Infoln(what)
}

// lerr is an ERROR log output with fields.
func lerr(msg string, err error, fields params) {
	if err == nil {
		return
	}
	log.WithError(errors.New(err.Error())).WithFields(logrus.Fields(fields)).Errorln(msg)
}

// l is an INFO log output.
func l(what interface{}) {
	log.Infoln(what)
}
