package log

import "github.com/sirupsen/logrus"

func NewLogWithConfiguration() *logrus.Logger {
	log := logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(customFormatter)
	return log
}
