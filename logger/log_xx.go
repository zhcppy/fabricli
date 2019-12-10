/*
@Time 2019-08-30 09:37
@Author ZH

*/
package logger

//import (
//	"os"
//	"strconv"
//
//	"github.com/op/go-logging"
//)
//
//const loggerName = ""
//
//var log *logging.Logger
//
//func L() *logging.Logger {
//	return log
//}
//
//func init() {
//	level, _ := strconv.ParseInt(os.Getenv("LEVEL"), 10, 8)
//	var _format = logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05.000} ▶ %{level:.4s} %{id:03d}%{color:reset} %{message}`)
//	var _backend = logging.NewLogBackend(os.Stderr, "", 0)
//	var leveledBackend = logging.MultiLogger(logging.NewBackendFormatter(_backend, _format))
//	leveledBackend.SetLevel(logging.Level(level), loggerName)
//
//	log = logging.MustGetLogger(loggerName)
//	log.SetBackend(leveledBackend)
//}
//
//func SetLevel(level int) {
//	var _format = logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05.000} ▶ %{level:.4s} %{id:03d}%{color:reset} %{message}`)
//	var _backend = logging.NewLogBackend(os.Stderr, "", 0)
//	var leveledBackend = logging.MultiLogger(logging.NewBackendFormatter(_backend, _format))
//	leveledBackend.SetLevel(logging.Level(level), loggerName)
//
//	log = logging.MustGetLogger(loggerName)
//	log.SetBackend(leveledBackend)
//}
