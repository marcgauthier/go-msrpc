package credential

import (
    "fmt"

    "github.com/oiweiwei/go-msrpc/internal/krb5compat"
    "github.com/oiweiwei/go-msrpc/internal/krb5v9/credentials"
)

// CCacheV8 is retained as a legacy credential marker. The typed v8 CCache()
// return signature cannot be retained without importing the old v8 module;
// legacy credentials are recognized by their CCache() method at runtime.
type CCacheV8 interface { Credential }

type CCache interface {
    Credential
    CCache() *credentials.CCache
}

type ccacheCred struct {
    userCred
    ccache *credentials.CCache
    ccacheErr error
}

func (cred *ccacheCred) CCache() *credentials.CCache {
    if cred != nil { return cred.ccache }
    return nil
}

func (cred *ccacheCred) IsEmpty() bool {
    return cred == nil || cred.ccache == nil || len(cred.ccache.Credentials) == 0
}

func (cred *ccacheCred) Validate() error {
    if cred != nil && cred.ccacheErr != nil { return cred.ccacheErr }
    return nil
}

func NewFromCCache(un string, ccache any, opts ...Option) CCache {
    switch ccache := ccache.(type) {
    case *credentials.CCache:
        return NewFromCCacheV9(un, ccache, opts...)
    default:
        return NewFromCCacheV8(un, ccache, opts...)
    }
}

func NewFromCCacheV9(un string, ccache *credentials.CCache, opts ...Option) CCache {
    return &ccacheCred{userCred: parseUser(un, opts...), ccache: ccache}
}

// NewFromCCacheV8 converts a legacy gokrb5/v8 credentials cache without
// referencing that module. Existing callers may pass their v8 cache directly.
func NewFromCCacheV8(un string, oldCache any, opts ...Option) CCache {
    if oldCache == nil {
        return &ccacheCred{ccacheErr: fmt.Errorf("nil Kerberos credentials cache")}
    }
    cc := new(credentials.CCache)
    if err := krb5compat.Copy(cc, oldCache); err != nil {
        return &ccacheCred{ccacheErr: fmt.Errorf("convert legacy credentials cache: %w", err)}
    }
    return NewFromCCacheV9(un, cc, opts...)
}
