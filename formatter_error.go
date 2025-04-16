package slogo

import (
	"log/slog"
	"reflect"
	"regexp"
	"runtime"

	slogformatter "github.com/samber/slog-formatter"
)

type ErrorStruct struct {
	Message    string
	Type       string
	Stacktrace string
}

func FormatError() slogformatter.Formatter {
	return slogformatter.FormatByType(func(err error) slog.Value {
		return anyValue(
			ErrorStruct{
				Message:    err.Error(),
				Type:       reflect.TypeOf(err).String(),
				Stacktrace: stacktrace(),
			},
		)
	})
}

var reStacktrace = regexp.MustCompile(`log/slog.*\n`)

func stacktrace() string {
	stackInfo := make([]byte, 1024*1024)

	if stackSize := runtime.Stack(stackInfo, false); stackSize > 0 {
		traceLines := reStacktrace.Split(string(stackInfo[:stackSize]), -1)
		if len(traceLines) > 0 {
			return traceLines[len(traceLines)-1]
		}
	}

	return ""
}
