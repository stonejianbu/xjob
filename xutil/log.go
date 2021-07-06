package xutil

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

var LogLevel = log.DebugLevel
var LogOut = os.Stdout
var LogOnce sync.Once
var JLog *log.Logger

func Log() *log.Logger {
	LogOnce.Do(func() {
		logger := log.New()
		logger.SetFormatter(&log.JSONFormatter{})
		logger.SetOutput(LogOut)
		logger.SetLevel(LogLevel)
		JLog = logger
	})
	return JLog
}
