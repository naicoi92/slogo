package slogo

import (
	"io"
	"log/slog"
	"time"

	slogformatter "github.com/samber/slog-formatter"
	slogmulti "github.com/samber/slog-multi"
	"gitlab.com/greyxor/slogor"
)

type SlogoHandler struct {
	options    *slog.HandlerOptions
	writer     io.Writer
	handlers   []slog.Handler
	formatters []slogformatter.Formatter
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions, options ...option) slog.Handler {
	h := &SlogoHandler{
		writer:     w,
		options:    opts,
		handlers:   []slog.Handler{},
		formatters: []slogformatter.Formatter{},
	}
	for _, opt := range options {
		opt(h)
	}
	formatters := h.getFormatters()
	handler := h.getHandler()
	return slogformatter.NewFormatterHandler(formatters...)(handler)
}

func (h *SlogoHandler) getHandler() slog.Handler {
	if len(h.handlers) == 0 {
		return slogor.NewHandler(
			h.writer,
			slogor.SetLevel(slog.LevelInfo),
			slogor.SetTimeFormat(time.Stamp),
			slogor.ShowSource())
	}
	if len(h.handlers) == 1 {
		return h.handlers[0]
	}
	return slogmulti.Fanout(h.handlers...)
}

func (h *SlogoHandler) getFormatters() []slogformatter.Formatter {
	if len(h.formatters) == 0 {
		return []slogformatter.Formatter{
			FormatStruct(),
			FormatError(),
		}
	}
	return h.formatters
}
