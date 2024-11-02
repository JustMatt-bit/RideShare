package core

import (
	"github.com/sirupsen/logrus"
)

type CtxKey string

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	CtxLog  CtxKey = "logger"
	CtxAuth CtxKey = "auth"
	CtxDB   CtxKey = "db"
)

func SetupLogging() *logrus.Entry {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
	return logrus.NewEntry(log)
}
