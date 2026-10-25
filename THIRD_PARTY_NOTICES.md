# Third-party code attribution

This fork copies the upstream implementations listed below into its
`internal/` directory to reduce external Go module dependencies.

| Internal directory | Upstream release | License | Changes |
| --- | --- | --- | --- |
| `internal/dns` | [miekg/dns v1.1.73](https://github.com/miekg/dns/tree/v1.1.73) | BSD 3-Clause | Moved source tree; callers changed to internal import path |
| `internal/smbfork` | [oiweiwei/go-smb2.fork v1.0.2](https://github.com/oiweiwei/go-smb2.fork/tree/v1.0.2) | BSD-style | Rewrote imports to nested internal packages and BER implementation; excluded generator/tests |
| `internal/ber` | [geoffgarside/ber v1.1.0](https://github.com/geoffgarside/ber/tree/v1.1.0) | BSD 3-Clause | Moved source under internal/ber |
| `internal/oem` | [oiweiwei/go-oem v1.0.0](https://github.com/oiweiwei/go-oem/tree/v1.0.0) | MIT | Rewrote EBCDIC imports and nested code-page package; excluded table generator |
| `internal/ebcdic` | [indece-official/go-ebcdic v1.2.0](https://github.com/indece-official/go-ebcdic/tree/v1.2.0) | MIT | Moved encoding code/data to internal/ebcdic |

All corresponding upstream `LICENSE` files are preserved in those directories.
Additional copyrights and license headers found in individual source files remain
intact. This notice is supplementary to the original licenses.

The imported DNS and SMB runtimes are initially kept broadly intact because
their protocol decoders and transport components reference many internal types.
Removing unreferenced functionality should follow successful compile, protocol
regression and cross-platform testing, not precede it.

## External module dependencies remaining

The permitted upstream `github.com/oiweiwei/gokrb5.fork/v9` imports a number
of third-party libraries. Source migration of DNS, SMB and OEM does **not**
remove those requirements. They must be handled through a separately audited
Kerberos refactor, retaining licenses and validating SPNEGO, ticket
cryptography and authentication interoperability.

## Public Go API change

The exported DNS record conversion APIs now use the internal DNS `RR` types
rather than the original `github.com/miekg/dns` types. Logger APIs were also
changed in the previous migration. These changes can break applications
which reference the previous external types; consumers require a migration
guide, or compatibility wrappers where practical.


## Kerberos and remaining utility source

| Internal directory | Source release | License | Changes |
| --- | --- | --- | --- |
| `internal/krb5v9` | [oiweiwei/gokrb5.fork v9.0.8](https://github.com/oiweiwei/gokrb5.fork/tree/v9.0.8/v9) | Apache-2.0 | Rehomed Kerberos v9 sources and rewrote imports; retained `LICENSE` and `NOTICE` |
| `internal/krbcompat/asn1`, `internal/krbcompat/pbkdf2` | [jcmturner/gofork v1.7.6](https://github.com/jcmturner/gofork/tree/v1.7.6) | BSD-style | Copied the ASN.1 and PBKDF2 implementations only; common license in `internal/krbcompat/LICENSE` |
| `internal/krbcompat/aescts` | [jcmturner/aescts v2.0.0](https://github.com/jcmturner/aescts/tree/v2.0.0/v2) | Apache-2.0 | Kept AES-CTS implementation and adapted RFC 3962 test vectors to Go `testing` |
| `internal/krbcompat/dnsutils` | [jcmturner/dnsutils v2.0.0](https://github.com/jcmturner/dnsutils/tree/v2.0.0/v2) | Apache-2.0 | Kept SRV/DNS configuration helper |
| `internal/krbcompat/goidentity` | [jcmturner/goidentity v6.0.1](https://github.com/jcmturner/goidentity/tree/v6.0.1/v6) | Apache-2.0 | Replaced HashiCorp UUID with internal `crypto/rand` implementation |
| `internal/krbcompat/rpc` | [jcmturner/rpc v2.0.3](https://github.com/jcmturner/rpc/tree/v2.0.3/v2) | Apache-2.0 | Kept Microsoft NDR and PAC types, with internal imports |
| `internal/floatmath` | [oiweiwei/go-math v1.0.0](https://github.com/oiweiwei/go-math/tree/v1.0.0) | MIT | Moved floating-point conversions used by NDR |

Each copied upstream codebase retains its complete license file. The Kerberos
Apache-2.0 NOTICE file is retained as `internal/krb5v9/NOTICE`. The new
`internal/krbcompat/uuid` helper is independently implemented with the Go
standard library and does not include code from HashiCorp's MPL-2.0 module.

Import-path modifications and test adaptations are changes to these sources;
this attribution does not imply that their original authors endorse this fork.
