package rex

import "github.com/Sirupsen/logrus"

func Info(message string) {
	logrus.Info(message)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Debug(message string) {
	logrus.Debug(message)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Warn(message string) {
	logrus.Warn(message)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(message string) {
	logrus.Error(message)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(message string) {
	logrus.Fatal(message)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Panic(message string) {
	logrus.Panic(message)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}
