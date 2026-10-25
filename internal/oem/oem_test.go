package oem

import (
    "bytes"
    "testing"
)

func TestCP437RoundTrip(t *testing.T) {
    codec := CP437_DOSLatinUS()
    encoded, err := codec.Encode("MURMUR")
    if err != nil { t.Fatal(err) }
    if !bytes.Equal(encoded, []byte("MURMUR")) { t.Fatalf("unexpected encoding %v", encoded) }
    decoded, err := codec.Decode(encoded)
    if err != nil { t.Fatal(err) }
    if decoded != "MURMUR" { t.Fatalf("round trip: %q", decoded) }
}

func TestEBCDIC37RoundTrip(t *testing.T) {
    codec := EBCDIC037()
    encoded, err := codec.Encode("HELLO")
    if err != nil { t.Fatal(err) }
    if len(encoded) != 5 || bytes.Equal(encoded, []byte("HELLO")) {
        t.Fatalf("unexpected EBCDIC encoded data %x", encoded)
    }
    decoded, err := codec.Decode(encoded)
    if err != nil { t.Fatal(err) }
    if decoded != "HELLO" { t.Fatalf("round trip: %q", decoded) }
}
