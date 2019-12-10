/*
@Time 2019-08-30 11:02
@Author ZH

*/
package logger

import (
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerName = os.Getenv("PROJECT_NAME")

var log *zap.SugaredLogger

func L() *zap.SugaredLogger {
	return log
}

func init() {
	level, _ := strconv.ParseInt(os.Getenv("LEVEL"), 10, 8)
	logConfig := zap.NewDevelopmentConfig()
	logConfig.DisableStacktrace = true
	logConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log = logger.Named(loggerName).WithOptions().Sugar()
}

func SetLevel(level int) {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.DisableStacktrace = true
	logConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log = logger.Named(loggerName).WithOptions().Sugar()
}
