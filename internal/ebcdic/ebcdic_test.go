package ebcdic

import "testing"

func TestCodePage500RoundTrip(t *testing.T) {
    encoded, err := Encode("ABC", EBCDIC500)
    if err != nil { t.Fatal(err) }
    decoded, err := Decode(encoded, EBCDIC500)
    if err != nil { t.Fatal(err) }
    if decoded != "ABC" { t.Fatalf("decoded %q, want ABC", decoded) }
}

func TestUnknownCodePage(t *testing.T) {
    if _, err := Encode("abc", -1); err == nil {
        t.Fatal("expected unsupported code page to fail")
    }
}
