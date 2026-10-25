package htmlquery

import (
    "strings"
    "testing"
)

func TestSelectors(t *testing.T) {
    const page = "<html><div class='content'><p>Intro</p><ul><li><a href='one'>Alpha</a></li></ul></div><table><tbody><tr><td><p>Row</p></td></tr></tbody></table></html>"
    doc, err := NewDocumentFromReader(strings.NewReader(page))
    if err != nil { t.Fatal(err) }
    testCases := []struct{ css, want string }{
        {"div.content > p", "Intro"},
        {"div.content ul li a", "Alpha"},
        {"table tbody tr", "Row"},
        {"td > p", "Row"},
        {"table tr", "Row"},
    }
    for _, tc := range testCases {
        got := doc.Find(tc.css).Text()
        if got != tc.want { t.Errorf("%s: got %q, want %q", tc.css, got, tc.want) }
    }
    a := doc.Find("div.content ul li a")
    if href, ok := a.Attr("href"); !ok || href != "one" {
        t.Fatalf("href = %q, %v", href, ok)
    }
    count := 0
    doc.Find("table tr").Each(func(_ int, _ *Selection) { count++ })
    if count != 1 { t.Fatalf("expected 1 row, got %d", count) }
}
