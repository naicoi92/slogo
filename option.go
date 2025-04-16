package slogo

import (
	"log/slog"
	"time"

	slogformatter "github.com/samber/slog-formatter"
	"gitlab.com/greyxor/slogor"
)

type Option func(*SlogoHandler)

func WithSlogor() Option {
	return func(h *SlogoHandler) {
		options := []slogor.OptionFn{
			slogor.SetTimeFormat(time.Stamp),
		}
		if h.options.AddSource {
			options = append(options, slogor.ShowSource())
		}
		if h.options.Level != nil {
			options = append(options, slogor.SetLevel(h.options.Level))
		}
		h.handlers = append(h.handlers, slogor.NewHandler(h.writer, options...))
	}
}

func WithSlogHandler(handler slog.Handler) Option {
	return func(h *SlogoHandler) {
		h.handlers = append(h.handlers, handler)
	}
}

func WithFormatter(formatter slogformatter.Formatter) Option {
	return func(h *SlogoHandler) {
		h.formatters = append(h.formatters, formatter)
	}
}
