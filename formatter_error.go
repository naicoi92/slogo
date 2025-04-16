package slogo

import (
	"log/slog"
	"reflect"

	slogformatter "github.com/samber/slog-formatter"
)

type ErrorStruct struct {
	Type    string
	Message string
}

func FormatError() slogformatter.Formatter {
	return slogformatter.FormatByType[error](func(err error) slog.Value {
		e := ErrorStruct{
			Type:    reflect.TypeOf(err).String(),
			Message: err.Error(),
		}
		return anyValue(e)
	})
}
