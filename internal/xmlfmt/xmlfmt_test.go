package xmlfmt

import (
    "strings"
    "testing"
)

func TestFormatXML(t *testing.T) {
    result := FormatXML("<root><child>value</child></root>", "", "  ")
    if !strings.Contains(result, "  <child>value</child>") {
        t.Fatalf("unexpected format %q", result)
    }
}
func TestInvalidXML(t *testing.T) {
    src := "<not-valid"
    if got := FormatXML(src, "", "  "); got != src {
        t.Fatalf("expected original invalid XML, got %q", got)
    }
}
