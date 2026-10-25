package wellknownyaml

import (
    "os"
    "testing"
)

func TestUnmarshal(t *testing.T) {
    input := "- uuid: 00000000-0000-0000-0000-000000000000\n" +
        "  domain: MC-CCFG\n  annotation: 'Server''s config'\n  if_name: Interface\n" +
        "  if_versions: [v1.0]\n  endpoints:\n  - '\\\\pipe\\\\endpoint'\n" +
        "  keywords:\n  - x\n- uuid: abc\n  endpoints: []\n"
    var got []Entry
    if err := Unmarshal([]byte(input), &got); err != nil { t.Fatal(err) }
    if len(got) != 2 || got[0].Annotation != "Server's config" || len(got[0].Endpoints) != 1 || len(got[1].Endpoints) != 0 {
        t.Fatalf("unexpected parsed entries: %#v", got)
    }
}

func TestRepositoryEndpointData(t *testing.T) {
    input, err := os.ReadFile("../../msrpc/well_known/data/well_known_endpoints.yaml")
    if err != nil { t.Fatal(err) }
    var got []Entry
    if err := Unmarshal(input, &got); err != nil { t.Fatal(err) }
    if len(got) != 784 {
        t.Fatalf("got %d well-known endpoints, want 784", len(got))
    }
    if got[0].UUID != "00000000-0000-0000-0000-000000000000" {
        t.Fatalf("unexpected first UUID %q", got[0].UUID)
    }
}
