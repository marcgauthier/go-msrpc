package record

import (
    "testing"

    "github.com/oiweiwei/go-msrpc/internal/dns"
)

// These cases form a regression baseline for extracting the DNS RR types
// from the third-party DNS library into an internal implementation.
func TestDNSRecordRoundTrip(t *testing.T) {
    cases := []string{
        ". 300 IN A 192.0.2.1",
        ". 300 IN AAAA 2001:db8::1",
        ". 300 IN NS ns.example.",
        ". 300 IN CNAME host.example.",
        ". 300 IN SOA ns.example. admin.example. 2026090101 3600 600 604800 60",
        ". 300 IN MX 10 mail.example.",
        ". 300 IN TXT \"hello\"",
        ". 300 IN SRV 1 2 443 service.example.",
        ". 300 IN PTR host.example.",
    }
    for _, line := range cases {
        t.Run(line, func(t *testing.T) {
            original, err := dns.NewRR(line)
            if err != nil { t.Fatal(err) }
            rec, err := NewRecordFromRR(original)
            if err != nil { t.Fatal(err) }
            back := rec.RR()
            if back == nil {
                t.Fatal("converted record is nil")
            }
            if back.Header().Rrtype != original.Header().Rrtype {
                t.Fatalf("record type: got %d, want %d", back.Header().Rrtype, original.Header().Rrtype)
            }
            if back.Header().Ttl != original.Header().Ttl {
                t.Fatalf("TTL: got %d, want %d", back.Header().Ttl, original.Header().Ttl)
            }
        })
    }
}
