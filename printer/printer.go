/*
@Time 2019-09-21 17:07
@Author ZH

*/
package printer

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/zhcppy/fabricli/jsonp"
)

type Printer interface {
	Success(format string, a ...interface{})
	Info(format string, a ...interface{})
	Fail(format string, a ...interface{})
	Error(format string, a ...interface{})
	Warn(format string, a ...interface{})
}

func Success(format string, a ...interface{}) {
	color.Blue(format, a...)
}

func Info(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func Fail(format string, a ...interface{}) {
	color.Magenta(format, a...)
}

func Error(format string, a ...interface{}) {
	color.Red(format, a...)
}

func Warn(format string, a ...interface{}) {
	color.Yellow(format, a...)
}

func JSON(v interface{}) {
	bytes, err := jsonp.Marshal(v)
	if err != nil {
		Error(err.Error())
	}
	fmt.Println(string(bytes))
}

func JSONStr(v interface{}) string {
	bytes, err := jsonp.Marshal(v)
	if err != nil {
		Error(err.Error())
	}
	return string(bytes)
}

type ColorFormatter struct {
	color *color.Color
}

func NewColorFmt() *ColorFormatter {
	return &ColorFormatter{
		color: color.New(color.FgBlue, color.Bold),
	}
}

func (cf *ColorFormatter) Success(format string, a ...interface{}) {
}

func (cf *ColorFormatter) Info(format string, a ...interface{}) {
}

func (cf *ColorFormatter) Fail(format string, a ...interface{}) {
}

func (cf *ColorFormatter) Error(format string, a ...interface{}) {
}

func (cf *ColorFormatter) Warn(format string, a ...interface{}) {
}
