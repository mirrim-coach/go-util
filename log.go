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
		hostname, err := os.Hostname()

		if err != nil {
			hostname = "unnamed_app"
		}

		hook, err := lSyslog.NewSyslogHook("udp", GetEnvVariable("LOG_PATH", "logs7.papertrailapp.com:51074"), syslog.LOG_INFO, GetEnvVariable("HOSTNAME", hostname))

		if err == nil {
			log.Hooks.Add(hook)
		}

		logger = log

	}

	return logger
}
