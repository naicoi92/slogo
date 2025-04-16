package slogo

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	slogformatter "github.com/samber/slog-formatter"
)

func FormatStruct() slogformatter.Formatter {
	return slogformatter.FormatByKind(slog.KindAny, func(v slog.Value) slog.Value {
		return anyValue(v.Any())
	})
}

func anyValue(s any) slog.Value {
	v := reflect.ValueOf(s)
	if ok, result := callStringFn(v, "String", "ToString"); ok {
		return slog.StringValue(result)
	}
	t := reflect.TypeOf(s)
	switch t.Kind() {
	case reflect.Ptr:
		return anyValue(v.Elem())
	case reflect.Array:
		return arrayValue(v)
	case reflect.Slice:
		return arrayValue(v)
	case reflect.Struct:
		return structValue(v)
	default:
		return slog.StringValue(fmt.Sprintf("%v", v.Interface()))
	}
}

func structValue(v reflect.Value) slog.Value {
	t := v.Type()
	var values []slog.Value
	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)
		if field.PkgPath != "" {
			continue
		}
		switch field.Tag.Get("slog") {
		case "-":
			continue
		case "restrict":
			values = append(values, slog.StringValue(fmt.Sprintf("%s=[REDACTED]", field.Name)))
			continue
		}
		if isEmptyValue(value) {
			continue
		}
		values = append(
			values,
			slog.StringValue(fmt.Sprintf("%s=%v", field.Name, anyValue(value).String())),
		)
	}
	fields := []string{}
	for _, v := range values {
		fields = append(fields, v.String())
	}
	return slog.StringValue("[" + strings.Join(fields, ", ") + "]")
}

func arrayValue(v reflect.Value) slog.Value {
	var values []slog.Value
	for i := range v.Len() {
		values = append(values, structValue(v.Index(i)))
	}
	if len(values) == 0 {
		return slog.StringValue("[]")
	}
	fields := []string{}
	for _, v := range values {
		fields = append(fields, v.String())
	}
	return slog.StringValue("[" + strings.Join(fields, ", ") + "]")
}

func callStringFn(v reflect.Value, methods ...string) (bool, string) {
	for _, methodString := range methods {
		if method := v.MethodByName(methodString); method.IsValid() {
			if result := method.Call(nil); len(result) > 0 {
				return true, result[0].Interface().(string)
			}
		}
	}
	return false, ""
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return v.Len() == 0
	case reflect.Map, reflect.Chan, reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	default:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}
