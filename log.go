package goutils

import (
	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
	"os"
)

var (
	logger *logrus.Logger
)

//Logger returns a new logger
func Logger() *logrus.Logger {
	if logger == nil {

		log := logrus.New()

		logLevel, err := logrus.ParseLevel(util.GetEnvVariable("LOG_LEVEL", "info"))
		if err == nil {
			log.SetLevel(logLevel)
		}

		hostname, err := os.Hostname()

		if err != nil {
			hostname = "unnamed_app"
		}

		hook, err := lSyslog.NewSyslogHook("udp", GetEnvVariable("LOG_HOST", "logs7.papertrailapp.com:51074"), syslog.LOG_INFO, GetEnvVariable("HOSTNAME", hostname))

		if err == nil && GetEnvVariable("ENABLE_PAPERTRAIL", "false") == "true" {
			log.Hooks.Add(hook)
		}

		logger = log

	}

	return logger
}
