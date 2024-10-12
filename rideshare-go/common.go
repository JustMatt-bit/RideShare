package main

import (
	"github.com/sirupsen/logrus"
)

func setupLogging() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
	return log
}
