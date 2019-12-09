package logger

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

//func init() {
//	log.SetReportCaller(true)
//}

func SetReportCaller(reportCaller bool) {
	log.SetReportCaller(reportCaller)
}
func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	log.Infoln(args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func AddHook(hook logrus.Hook) {
	log.AddHook(hook)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}
