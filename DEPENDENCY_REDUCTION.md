# Dependency internalization audit

**Target:** No non-permitted Go modules outside the Go standard library, `golang.org/*`, `google.golang.org/*` and `github.com/oiweiwei/*`.

## Current declared module state

The `dependency-internalization` branch currently declares only four
`golang.org/x/*` modules in `go.mod`:

- `golang.org/x/crypto`
- `golang.org/x/text`
- `golang.org/x/net`
- `golang.org/x/sys` (indirect)

There are **zero non-permitted modules declared** in `go.mod` and
`go.sum` contains checksums only for permitted modules. This is a
source and manifest audit, **not yet a successful build or verified Go module
dependency graph**.

## Internal packages

| Package | Origin | Purpose |
| --- | --- | --- |
| `internal/logging` | Go `log/slog` | RPC structured logging without Zerolog |
| `internal/dns` | miekg/dns | DNS resource-record parsing and formatting |
| `internal/smbfork` | oiweiwei/go-smb2.fork | SMB2/3 transport and named pipes |
| `internal/ber` | geoffgarside/ber | SMB SPNEGO ASN.1 decoding |
| `internal/oem` | oiweiwei/go-oem | NTLM/Netlogon OEM encodings |
| `internal/ebcdic` | indece-official/go-ebcdic | IBM EBCDIC code pages |
| `internal/krb5v9` | oiweiwei/gokrb5.fork/v9 | Kerberos 5/SPNEGO implementation |
| `internal/krbcompat` | JCM Turner ASN.1, PBKDF2, AES-CTS, DNS, GoIdentity, PAC/NDR; Go `crypto/rand` | Only the supporting Kerberos components needed by the imported code |
| `internal/floatmath` | oiweiwei/go-math | NDR non-IEEE floats |
| `internal/krb5compat` | New adapter | Legacy v8-shaped CCache/keytab/config conversion |
| `internal/htmlquery` | New Go `x/net/html` adapter | HTML scraping in generators |
| `internal/xmlfmt` | New standard-library formatter | XML examples |
| `internal/wellknownyaml` | New fixed-schema parser | Well-known endpoint data generator |

See `THIRD_PARTY_NOTICES.md` for provenance, license requirements, copyright and attribution.

## Verification required before merge

1. Enable GitHub Actions on this fork if disabled. The branch CI workflow
   checks the **full** module graph and the allowlist, then runs `go test ./...`
   and `go vet ./...` with the Go toolchain from `go.mod` (currently Go 1.26).
2. Verify `go mod tidy` makes no unexpected changes, and run
   `go list -m all` and `go mod verify`.
3. Exercise Kerberos AS/TGS/AP exchanges, GSSAPI/SPNEGO MIC and wrap/unwrap,
   service authentication, credential-cache/keytab conversion, and PAC decoding.
4. Verify Windows SMB2/3 authenticated named pipes, DNS RR round trips,
   NTLM and Netlogon, and generated MSRPC stubs.
5. Check public API compatibility. Moving types from the external DNS, SMB,
   and logging packages into `internal/` **changes exported method and field
   types**, potentially breaking downstream applications. A migration guide or
   compatibility layer is required.
6. Review the copied code for pruning opportunities only after successful
   compilation and protocol integration testing.

**Status:** Draft, source-level work performed, not yet successfully compiled
or integration-tested in this environment. Do not merge into `main` yet.
