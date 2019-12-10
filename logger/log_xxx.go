/*
@Time 2019-08-30 10:35
@Author ZH

*/
package logger

//import (
//	"os"
//	"strconv"
//
//	"github.com/sirupsen/logrus"
//)
//
//const loggerName = ""
//
//var log *logrus.Logger
//
//func L() *logrus.Logger {
//	if log == nil {
//		log = logrus.New()
//	}
//	return log
//}
//func init() {
//	level, _ := strconv.ParseInt(os.Getenv("LEVEL"), 10, 8)
//	log := logrus.New()
//	log.SetLevel(logrus.Level(level))
//	log.SetFormatter(&logrus.TextFormatter{})
//	//log.SetFormatter(&logrus.JSONFormatter{})
//	log.SetReportCaller(true)
//}
//
//func SetLevel(level int) {
//	log := logrus.New()
//	log.SetLevel(logrus.Level(level))
//	log.SetFormatter(&logrus.TextFormatter{})
//	//log.SetFormatter(&logrus.JSONFormatter{})
//	log.SetReportCaller(true)
//}
