package slogo

import (
	"log/slog"
	"reflect"

	slogformatter "github.com/samber/slog-formatter"
)

type ErrorStruct struct {
	Message string
	Type    string
}

func FormatError() slogformatter.Formatter {
	return slogformatter.FormatByType(func(err error) slog.Value {
		e := ErrorStruct{
			Message: err.Error(),
			Type:    reflect.TypeOf(err).String(),
		}
		return anyValue(e)
	})
}
