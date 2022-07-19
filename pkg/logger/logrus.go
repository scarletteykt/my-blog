package logger

import (
	"github.com/sirupsen/logrus"
)

type Logrus struct {
	logger *logrus.Logger
}

func NewLogger() Logger {
	return &Logrus{logrus.New()}
}

func (l *Logrus) Debug(msg ...interface{}) {
	l.logger.Debug(msg...)
}

func (l *Logrus) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logrus) Info(msg ...interface{}) {
	l.logger.Info(msg...)
}

func (l *Logrus) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logrus) Warn(msg ...interface{}) {
	l.logger.Warn(msg...)
}

func (l *Logrus) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *Logrus) Error(msg ...interface{}) {
	l.logger.Error(msg...)
}

func (l *Logrus) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logrus) Fatal(msg ...interface{}) {
	l.logger.Fatal(msg...)
}

func (l *Logrus) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
