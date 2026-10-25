package credential

import (
    "testing"

    "github.com/oiweiwei/gokrb5.fork/v9/credentials"
    "github.com/oiweiwei/gokrb5.fork/v9/keytab"
)

type oldPrincipal struct { Realm string }
type oldCache struct {
    Version uint8
    DefaultPrincipal oldPrincipal
    Credentials []*oldCredential
    Path string
}
type oldCredential struct { Client oldPrincipal; Server oldPrincipal }
type oldKeytab struct { Entries []oldEntry }
type oldEntry struct { KVNO uint32 }
type legacyCacheCredential struct { Credential; cache *oldCache }
func (c legacyCacheCredential) CCache() *oldCache { return c.cache }
type legacyKeytabCredential struct { Credential; tab *oldKeytab }
func (c legacyKeytabCredential) Keytab() *oldKeytab { return c.tab }

func TestLegacyCCacheConversion(t *testing.T) {
    old := &oldCache{Version: 4, DefaultPrincipal: oldPrincipal{Realm: "EXAMPLE.COM"},
        Credentials: []*oldCredential{{Client: oldPrincipal{Realm: "EXAMPLE.COM"}}}}
    got := NewFromCCacheV8("user@EXAMPLE.COM", old)
    if err := got.(*ccacheCred).Validate(); err != nil { t.Fatal(err) }
    cache := got.CCache()
    if cache == nil || len(cache.Credentials) != 1 || cache.DefaultPrincipal.Realm != "EXAMPLE.COM" {
        t.Fatalf("legacy ccache not converted: %#v", cache)
    }
    upgraded := V8ToV9(legacyCacheCredential{Credential: New("user@EXAMPLE.COM"), cache: old})
    if _, ok := upgraded.(CCache); !ok { t.Fatalf("expected converted ccache, got %T", upgraded) }
    if original := NewFromCCacheV9("x", &credentials.CCache{}); V8ToV9(original) != original {
        t.Fatal("native v9 CCache should not be converted")
    }
}
func TestLegacyKeytabConversion(t *testing.T) {
    old := &oldKeytab{Entries: []oldEntry{{KVNO: 7}}}
    got := NewFromKeytabV8("user@EXAMPLE.COM", old)
    if err := got.(*keytabCred).Validate(); err != nil { t.Fatal(err) }
    if got.Keytab() == nil || len(got.Keytab().Entries) != 1 {
        t.Fatalf("legacy keytab not converted: %#v", got.Keytab())
    }
    upgraded := V8ToV9(legacyKeytabCredential{Credential: New("user@EXAMPLE.COM"), tab: old})
    if _, ok := upgraded.(Keytab); !ok { t.Fatalf("expected converted keytab, got %T", upgraded) }
    if original := NewFromKeytabV9("x", &keytab.Keytab{}); V8ToV9(original) != original {
        t.Fatal("native v9 keytab should not be converted")
    }
}
func TestInvalidLegacyData(t *testing.T) {
    if err := NewFromCCacheV8("x", nil).(*ccacheCred).Validate(); err == nil { t.Fatal("expected nil cache error") }
    if err := NewFromKeytabV8("x", 7).(*keytabCred).Validate(); err == nil { t.Fatal("expected keytab conversion error") }
}
