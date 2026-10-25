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
