package credential

import (
    "fmt"

    "github.com/oiweiwei/go-msrpc/internal/krb5compat"
    "github.com/oiweiwei/go-msrpc/internal/krb5v9/keytab"
)

// KeytabV8 is a legacy credential marker. Legacy Keytab() methods are
// detected reflectively without importing the gokrb5/v8 module.
type KeytabV8 interface { Credential }

// Keytab is a Kerberos 5 keytab credential.
type Keytab interface {
    Credential
    Keytab() *keytab.Keytab
}

type keytabCred struct {
    userCred
    keytab *keytab.Keytab
    keytabErr error
}

func (cred *keytabCred) Keytab() *keytab.Keytab {
    if cred != nil { return cred.keytab }
    return &keytab.Keytab{}
}

func (cred *keytabCred) Validate() error {
    if cred != nil && cred.keytabErr != nil { return cred.keytabErr }
    return nil
}

func (cred *keytabCred) IsEmpty() bool {
    return cred == nil || cred.keytab == nil || len(cred.keytab.Entries) == 0
}

func NewFromKeytabFile(un string, keytabFile string, opts ...Option) Keytab {
    kt, err := keytab.Load(keytabFile)
    if err != nil { return &keytabCred{keytabErr: err} }
    return NewFromKeytabV9(un, kt, opts...)
}

func NewFromKeytab(un string, ktab any, opts ...Option) Keytab {
    switch ktab := ktab.(type) {
    case *keytab.Keytab:
        return NewFromKeytabV9(un, ktab, opts...)
    default:
        return NewFromKeytabV8(un, ktab, opts...)
    }
}

func NewFromKeytabV9(un string, kt *keytab.Keytab, opts ...Option) Keytab {
    return &keytabCred{userCred: parseUser(un, opts...), keytab: kt}
}

// NewFromKeytabV8 converts a legacy v8 keytab to v9 without importing v8.
func NewFromKeytabV8(un string, oldKeytab any, opts ...Option) Keytab {
    if oldKeytab == nil {
        return &keytabCred{keytabErr: fmt.Errorf("nil Kerberos keytab")}
    }
    kt := new(keytab.Keytab)
    if err := krb5compat.Copy(kt, oldKeytab); err != nil {
        return &keytabCred{keytabErr: fmt.Errorf("convert legacy keytab: %w", err)}
    }
    return NewFromKeytabV9(un, kt, opts...)
}
