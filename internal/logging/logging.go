// Package logging provides the small structured-logging API required by go-msrpc.
// The implementation uses only the Go standard library.
package logging

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "time"
)

// Logger is a lightweight structured logger. Its zero value discards messages.
type Logger struct {
    logger *slog.Logger
}

// New constructs a JSON logger with debug messages enabled.
func New(w io.Writer) Logger {
    if w == nil {
        return Logger{}
    }
    return Logger{logger: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))}
}

// Context builds a logger with fixed attributes.
type Context struct {
    logger Logger
    fields []any
}

func (l Logger) With() Context { return Context{logger: l} }

func (c Context) Str(key, value string) Context {
    c.fields = append(append([]any(nil), c.fields...), key, value)
    return c
}

func (c Context) Logger() Logger {
    if c.logger.logger == nil {
        return Logger{}
    }
    return Logger{logger: c.logger.logger.With(c.fields...)}
}

// Event represents one structured log record.
type Event struct {
    logger *slog.Logger
    level slog.Level
    fields []any
}

// LogObjectMarshaler is used by RPC headers, PDUs, and transport settings.
type LogObjectMarshaler interface {
    MarshalZerologObject(*Event)
}

func (l Logger) event(level slog.Level) *Event {
    return &Event{logger: l.logger, level: level}
}

func (l Logger) Debug() *Event { return l.event(slog.LevelDebug) }
func (l Logger) Info() *Event  { return l.event(slog.LevelInfo) }
func (l Logger) Warn() *Event  { return l.event(slog.LevelWarn) }
func (l Logger) Error() *Event { return l.event(slog.LevelError) }
func (l Logger) Err(err error) *Event { return l.Error().Err(err) }

func (e *Event) field(key string, value any) *Event {
    if e != nil {
        e.fields = append(e.fields, key, value)
    }
    return e
}

func (e *Event) Str(key, value string) *Event { return e.field(key, value) }
func (e *Event) Int(key string, value int) *Event { return e.field(key, value) }
func (e *Event) Uint8(key string, value uint8) *Event { return e.field(key, value) }
func (e *Event) Uint16(key string, value uint16) *Event { return e.field(key, value) }
func (e *Event) Uint32(key string, value uint32) *Event { return e.field(key, value) }
func (e *Event) Bool(key string, value bool) *Event { return e.field(key, value) }
func (e *Event) Dur(key string, value time.Duration) *Event { return e.field(key, value.String()) }
func (e *Event) Interface(key string, value any) *Event { return e.field(key, value) }
func (e *Event) Stringer(key string, value fmt.Stringer) *Event {
    if value == nil {
        return e.field(key, nil)
    }
    return e.field(key, value.String())
}
func (e *Event) RawJSON(key string, value []byte) *Event {
    if json.Valid(value) {
        return e.field(key, json.RawMessage(value))
    }
    return e.field(key, string(value))
}
func (e *Event) Err(err error) *Event {
    if err != nil {
        return e.field("error", err.Error())
    }
    return e
}
func (e *Event) EmbedObject(value LogObjectMarshaler) *Event {
    if e != nil && value != nil {
        value.MarshalZerologObject(e)
    }
    return e
}
func (e *Event) Msg(message string) {
    if e == nil || e.logger == nil {
        return
    }
    e.logger.Log(context.Background(), e.level, message, e.fields...)
}
func (e *Event) Msgf(format string, args ...any) {
    if e != nil {
        e.Msg(fmt.Sprintf(format, args...))
    }
}
