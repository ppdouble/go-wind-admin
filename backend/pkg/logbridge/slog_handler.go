package logbridge

import (
	"context"
	"log/slog"
	"strings"

	bootstraplog "github.com/tx7do/kratos-bootstrap/logger"
)

// Handler forwards Kratos v3 slog records to the existing configured logger.
// This keeps request trace fields and deployment log output in one pipeline.
type Handler struct {
	logger bootstraplog.Logger
	attrs  []any
	groups []string
}

func NewHandler(logger bootstraplog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) Enabled(context.Context, slog.Level) bool { return true }

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	fields := append([]any(nil), h.attrs...)
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(&fields, strings.Join(h.groups, "."), attr)
		return true
	})
	switch {
	case record.Level >= slog.LevelError:
		h.logger.Error(ctx, record.Message, fields...)
	case record.Level >= slog.LevelWarn:
		h.logger.Warn(ctx, record.Message, fields...)
	case record.Level >= slog.LevelInfo:
		h.logger.Info(ctx, record.Message, fields...)
	default:
		h.logger.Debug(ctx, record.Message, fields...)
	}
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := &Handler{logger: h.logger, attrs: append([]any(nil), h.attrs...), groups: append([]string(nil), h.groups...)}
	for _, attr := range attrs {
		appendAttr(&next.attrs, strings.Join(next.groups, "."), attr)
	}
	return next
}

func (h *Handler) WithGroup(name string) slog.Handler {
	next := &Handler{logger: h.logger, attrs: append([]any(nil), h.attrs...), groups: append([]string(nil), h.groups...)}
	if name != "" {
		next.groups = append(next.groups, name)
	}
	return next
}

func appendAttr(fields *[]any, prefix string, attr slog.Attr) {
	if attr.Equal(slog.Attr{}) {
		return
	}
	key := attr.Key
	if prefix != "" {
		key = prefix + "." + key
	}
	value := attr.Value.Resolve()
	if value.Kind() == slog.KindGroup {
		for _, child := range value.Group() {
			appendAttr(fields, key, child)
		}
		return
	}
	*fields = append(*fields, key, value.Any())
}
