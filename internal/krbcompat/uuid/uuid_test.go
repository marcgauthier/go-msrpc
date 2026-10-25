package uuid

import (
    "testing"
    "strings"
)

func TestGenerateUUID(t *testing.T) {
    first, err := GenerateUUID()
    if err != nil { t.Fatal(err) }
    second, err := GenerateUUID()
    if err != nil { t.Fatal(err) }
    if len(first)!=36 || first[14]!='4' || !strings.ContainsAny(first[19:20],"89ab") {
        t.Fatalf("invalid UUID v4 %q",first)
    }
    if first==second { t.Fatal("two UUIDs should not normally be identical") }
}
