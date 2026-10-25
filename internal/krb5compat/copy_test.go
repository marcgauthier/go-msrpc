package krb5compat

import (
    "reflect"
    "testing"
    "time"
)

type oldName struct { NameType int32; NameString []string }
type newName struct { NameType int32; NameString []string }
type oldCred struct {
    Realm string
    PrincipalName oldName
    Keys []*oldKey
    When time.Time
}
type oldKey struct { KeyType int32; KeyValue []byte }
type newKey struct { KeyType int32; KeyValue []byte }
type newCred struct {
    Realm string
    PrincipalName newName
    Keys []*newKey
    When time.Time
    Optional int
}

func TestCopyLegacyStruct(t *testing.T) {
    now := time.Now()
    src := &oldCred{
        Realm: "EXAMPLE.COM", PrincipalName: oldName{1, []string{"service", "host"}},
        Keys: []*oldKey{{18, []byte{1, 2, 3}}}, When: now,
    }
    var got newCred
    if err := Copy(&got, src); err != nil { t.Fatal(err) }
    if got.Realm != src.Realm || !reflect.DeepEqual(got.PrincipalName.NameString, src.PrincipalName.NameString) ||
        got.Keys[0].KeyType != 18 || got.Keys[0].KeyValue[2] != 3 || !got.When.Equal(now) {
        t.Fatalf("conversion mismatch: %#v", got)
    }
}
func TestCopyInvalid(t *testing.T) {
    if err := Copy(nil, 5); err == nil { t.Fatal("expected error for nil destination") }
    var got newCred
    if err := Copy(&got, 5); err == nil { t.Fatal("expected incompatible-type error") }
}
