// Package logging exposes go-msrpc's minimal structured logger API.
// The implementation lives in internal/logging and has no external dependencies.
package logging

import (
    "io"
    internal "github.com/oiweiwei/go-msrpc/internal/logging"
)

type Logger = internal.Logger
type Event = internal.Event
type Context = internal.Context
type LogObjectMarshaler = internal.LogObjectMarshaler

func New(w io.Writer) Logger { return internal.New(w) }
