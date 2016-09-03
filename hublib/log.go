package hublib

import (
  "github.com/jdrivas/sl"
  "github.com/Sirupsen/logrus"
)

var (
  log = sl.New()
)

func init() {
  initLogs()
}

func SetLogLevel(l logrus.Level) {
  log.SetLevel(l)
}

func SetLogFormatter(f logrus.Formatter) {
  log.SetFormatter(f)
}

func initLogs() {
  formatter := new(sl.TextFormatter)
  formatter.FullTimestamp = true
  log.SetFormatter(formatter)
  // log.SetLevel(logrus.InfoLevel)
  log.SetLevel(logrus.DebugLevel)
}

