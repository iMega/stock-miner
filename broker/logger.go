package broker

import "github.com/sirupsen/logrus"

type logger struct {
	log logrus.FieldLogger
}

func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Infof(msg, keysAndValues...)
}

func (l *logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Errorf(msg, keysAndValues...)
}
