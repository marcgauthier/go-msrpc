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

## Outstanding

- **Direct external runtime dependency:** `github.com/miekg/dns`. Windows DNS record conversion in `msrpc/dnsp` currently handles dozens of RR types and exports `dns.RR` in public signatures. Migrating requires API review and round-trip tests for each supported record type.
- **Transitive dependencies of permitted modules:** These do not disappear when a direct `require` is removed. `github.com/oiweiwei/gokrb5.fork/v9` imports jcmturner, hashicorp and others; `github.com/oiweiwei/go-smb2.fork` uses `github.com/geoffgarside/ber`; `github.com/oiweiwei/go-oem` uses `github.com/indece-official/go-ebcdic`. Fully eliminating external downloads will require modifications to those implementations or carefully scoped module replacements.
- **API compatibility:** `WithLogger` now expects this fork's `logging.Logger`, not `zerolog.Logger`. Legacy v8 credential types are accepted through `any`/reflection; their original v8-specific static interfaces cannot be identical without retaining v8 imports.
- **Verification:** CI has been expanded but has not run in the fork. Go 1.26 build, tests, remote Windows RPC, SMB, NTLM and Kerberos interoperability are not yet verified.

## Completion criteria

1. No non-permitted direct imports or transitive modules required for the selected build.
2. `go mod tidy`, `go test ./...` and race tests pass with Go 1.26.
3. DNS RR conversion tests cover all previously supported record types.
4. Kerberos v9 and legacy v8 adapter fixtures, NTLM and remote RPC interop work.
5. Licenses and attribution are retained for any copied external implementation.
