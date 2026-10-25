package credential

import (
    "reflect"
    "strings"
)

// Generic credential.
type Credential interface {
	// User name.
	UserName() string
	// Domain name.
	DomainName() string
	// Workstation.
	Workstation() string
}

func parseUser(un string, opts ...Option) userCred {
	dn, un, wkst := parseDomainUserWorkstation(un, opts...)
	return userCred{
		userName:    un,
		domainName:  dn,
		workstation: wkst,
	}
}

// parseDomainUserWorkstation parses the username and variadic workstation argument and outputs
// the domain name, user name, and workstation.
func parseDomainUserWorkstation(un string, opts ...Option) (string, string, string) {

	wkst, dn := "", ""

	for _, opt := range opts {
		switch v := opt.(type) {
		case wkstOpt:
			wkst = string(v)
		case domainOpt:
			dn = string(v)
		}
	}

	// down-level logon name.
	if strings.Contains(un, "\\") {
		un := strings.SplitN(un, "\\", 2)
		if dn != "" {
			return dn, un[1], wkst
		}
		return un[0], un[1], wkst
	}

	if strings.Contains(un, "@") {
		un := strings.SplitN(un, "@", 2)
		if dn != "" {
			return dn, un[0], wkst
		}
		return un[1], un[0], wkst
	}

	return dn, un, wkst
}

// DomainName function returns the domain name from the user name.
func DomainName(un string) string {
	dn, _, _ := parseDomainUserWorkstation(un)
	return dn
}

func IsAnonymous(cred Credential) bool {
	if cred == nil {
		return true
	}
	p, ok := (any)(cred).(*passwordCred)
	if !ok {
		return false
	}
	return p.allowEmpty && p.userName == ""
}

// Anonymous function returns the anonymous password credentials.
func Anonymous() Password {
	return &passwordCred{allowEmpty: true}
}

func New(un string, opts ...Option) Credential {
	dn, un, wkst := parseDomainUserWorkstation(un, opts...)
	return &userCred{
		userName:    un,
		domainName:  dn,
		workstation: wkst,
	}
}

func NewFromString(s string, opts ...Option) Password {
	ss := strings.Split(s, "%")
	if len(ss) > 1 {
		return NewFromPassword(ss[0], ss[1], opts...)
	}
	return NewFromPassword(s, "", append(opts, AllowEmptyPassword())...)
}

// V8ToV9 upgrades a legacy credential implementing CCache() or Keytab().
// Reflection allows the old method's return type to remain *gokrb5/v8/...
// in calling applications without requiring a v8 import in this module.
func V8ToV9(cred Credential) Credential {
    if cred == nil { return nil }
    if _, ok := cred.(CCache); ok { return cred }
    if _, ok := cred.(Keytab); ok { return cred }

    value := reflect.ValueOf(cred)
    if value.Kind() == reflect.Pointer && value.IsNil() { return cred }

    for _, kind := range []string{"Keytab", "CCache"} {
        method := value.MethodByName(kind)
        if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
            continue
        }
        result := method.Call(nil)[0]
        if !result.IsValid() { continue }
        if (result.Kind() == reflect.Pointer || result.Kind() == reflect.Interface) && result.IsNil() {
            continue
        }
        opts := []Option{Workstation(cred.Workstation()), Domain(cred.DomainName())}
        if kind == "Keytab" {
            return NewFromKeytabV8(cred.UserName(), result.Interface(), opts...)
        }
        return NewFromCCacheV8(cred.UserName(), result.Interface(), opts...)
    }
    return cred
}
