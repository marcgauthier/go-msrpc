# Dependency reduction status

This fork aims to build with dependencies only from Go (standard library and `golang.org`), `google.golang.org`, and `github.com/oiweiwei`.

## Internal implementations

| Path | Replaces | Scope |
| --- | --- | --- |
| `internal/logging` | `github.com/rs/zerolog` | Structured JSON logging using `log/slog` |
| `internal/krb5compat` | direct `github.com/jcmturner/gokrb5/v8` imports | Legacy CCache, keytab, and Krb5 config conversion into allowed v9 types |
| `internal/htmlquery` | `github.com/PuerkitoBio/goquery` | Selected CSS selectors used by three generators; uses allowed `golang.org/x/net/html` |
| `internal/xmlfmt` | `github.com/go-xmlfmt/xmlfmt` | XML indentation for Event Log example programs |
| `internal/wellknownyaml` | generator's direct `gopkg.in/yaml.v3` use | Fixed-format YAML endpoint metadata only, not a general YAML parser |
| Standard library `testing` | direct Testify assertions | DRSUAPI test assertions |

All changes are under `dependency-internalization`. Source links using `github.com/oiweiwei/go-msrpc` remain because this fork retains the original module path.

## Internalized runtime sources (additional changes)

- `internal/dns`: miekg/dns v1.1.73, BSD-3-Clause. Windows DNS record conversions now use internal DNS types instead of `miekg/dns`.
- `internal/ber`: geoffgarside/ber v1.1.0, BSD-3-Clause.
- `internal/smbfork`: oiweiwei/go-smb2.fork v1.0.2, BSD-style; its nested imports point to internal package paths. RPC named pipes now use this implementation.
- `internal/oem`: oiweiwei/go-oem v1.0.0, MIT; NTLM and Netlogon now use this internal implementation.
- `internal/ebcdic`: indece-official/go-ebcdic v1.2.0, MIT.

Upstream licenses and attribution are maintained in the internal directories and summarized in `THIRD_PARTY_NOTICES.md`.

## Outstanding

- **Nine non-permitted indirect modules remain declared** in `go.mod`, predominantly from `github.com/oiweiwei/gokrb5.fork/v9`. None are declared as direct dependencies anymore. The complete Go module graph and build outputs must still be checked using `go list -m all` and `go mod tidy`.
- **Validation and migration:** Run `go test ./...` with Go 1.26, then SMB2/3, DNS, NTLM, SPNEGO and Kerberos interoperability against Windows. The new internal DNS and SMB types and the logger API may break importing programs; test/migrate those call sites.
- **Code pruning:** The copied DNS and SMB trees still include upstream implementation pieces that may be unused by the RPC library. Prune only after successful build and automated compatibility coverage.
- **GitHub Actions:** No GitHub Actions workflow run has been observed for the fork. The CI definition now triggers on this branch as well, when GitHub Actions is enabled.

**Status: work in progress; this branch is not verified production-ready.**

## Completion criteria

1. No non-permitted direct imports or transitive modules required for the selected build.
2. `go mod tidy`, `go test ./...` and race tests pass with Go 1.26.
3. DNS RR conversion tests cover all previously supported record types.
4. Kerberos v9 and legacy v8 adapter fixtures, NTLM and remote RPC interop work.
5. Licenses and attribution are retained for any copied external implementation.
