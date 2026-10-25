package wellknownyaml

import "testing"

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
