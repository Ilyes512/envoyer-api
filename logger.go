package envoyerapi

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type redactHandler struct {
	inner      slog.Handler
	redactKeys map[string]struct{}
}

func NewRedactHandler(inner slog.Handler, keys ...string) slog.Handler {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[http.CanonicalHeaderKey(k)] = struct{}{}
	}
	return &redactHandler{inner: inner, redactKeys: m}
}

func (h *redactHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *redactHandler) Handle(ctx context.Context, r slog.Record) error {
	r2 := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		r2.AddAttrs(h.rewrite(a))
		return true
	})

	return h.inner.Handle(ctx, r2)
}

func (h *redactHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &redactHandler{inner: h.inner.WithAttrs(attrs), redactKeys: h.redactKeys}
}

func (h *redactHandler) WithGroup(name string) slog.Handler {
	return &redactHandler{inner: h.inner.WithGroup(name), redactKeys: h.redactKeys}
}

func (h *redactHandler) rewrite(a slog.Attr) slog.Attr {
	if strings.Contains(a.Key, "header") && a.Value.Kind() == slog.KindAny {
		if hdrs, ok := a.Value.Any().(http.Header); ok {
			cl := hdrs.Clone()
			for k := range cl {
				if _, hit := h.redactKeys[http.CanonicalHeaderKey(k)]; hit {
					cl.Set(k, "<*REDACTED*>")
				}
			}
			return slog.Any(a.Key, cl)
		}
	}
	return a
}
