package logging

import (
    "bytes"
    "encoding/json"
    "errors"
    "testing"
)

type testObject struct{}
func (testObject) MarshalZerologObject(e *Event) { e.Uint32("call_id", 42) }

func TestStructuredLog(t *testing.T) {
    var b bytes.Buffer
    logger := New(&b).With().Str("source", "rpc").Logger()
    logger.Debug().EmbedObject(testObject{}).Err(errors.New("network")).Str("target", "host").Msgf("call %d", 42)
    var record map[string]any
    if err := json.Unmarshal(b.Bytes(), &record); err != nil {
        t.Fatal(err)
    }
    if record["source"] != "rpc" || record["call_id"] != float64(42) ||
        record["error"] != "network" || record["target"] != "host" ||
        record["message"] != "call 42" {
        t.Fatalf("unexpected log record: %#v", record)
    }
}
func TestZeroLogger(t *testing.T) {
    var logger Logger
    logger.Debug().Str("key", "value").Msg("ignored")
    logger.Err(nil).Msg("ignored")
}
