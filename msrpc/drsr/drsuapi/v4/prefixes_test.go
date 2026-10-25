package drsuapi

import (
    "encoding/asn1"
    "testing"
)

var TestData = []asn1.ObjectIdentifier{
    {1, 2, 840, 113556, 1, 4, 94},
    {2, 5, 6, 2},
    {1, 2, 840, 113556, 1, 2, 1},
    {1, 2, 840, 113556, 1, 3, 223},
    {1, 2, 840, 113556, 1, 5, 7000, 53},
}

func TestPrefixTable(t *testing.T) {
    prefixes := make(PrefixTable, 0)
    attrs := make([]uint32, len(TestData))

    for i, oid := range TestData {
        attr, err := prefixes.OIDToAttribute(oid)
        if err != nil {
            t.Fatalf("convert OID %v to attribute: %v", oid, err)
        }
        attrs[i] = attr
    }

    for i, want := range TestData {
        got, err := prefixes.AttributeToOID(attrs[i])
        if err != nil {
            t.Fatalf("convert attribute %d to OID: %v", attrs[i], err)
        }
        if !got.Equal(want) {
            t.Errorf("attribute %d: OID = %v, want %v", attrs[i], got, want)
        }
    }
}
