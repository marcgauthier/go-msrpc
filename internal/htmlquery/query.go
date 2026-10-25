// Package htmlquery implements the small HTML selection API needed by
// the Microsoft documentation code generators. No CSS engine is required:
// the supported selectors are tag names, tag.class, descendants and > children.
package htmlquery

import (
    "fmt"
    "io"
    "net/http"
    "strings"

    "golang.org/x/net/html"
)

type Selection struct { nodes []*html.Node }
type Document struct { *Selection }

// NewDocumentFromReader parses an HTML document from r.
func NewDocumentFromReader(r io.Reader) (*Document, error) {
    n, err := html.Parse(r)
    if err != nil { return nil, err }
    return &Document{Selection: &Selection{nodes: []*html.Node{n}}}, nil
}

// NewDocument downloads and parses an HTML document.
func NewDocument(url string) (*Document, error) {
    resp, err := http.Get(url)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("get %s: %s", url, resp.Status)
    }
    return NewDocumentFromReader(resp.Body)
}

func matches(n *html.Node, selector string) bool {
    if n.Type != html.ElementNode { return false }
    parts := strings.SplitN(selector, ".", 2)
    if parts[0] != "" && !strings.EqualFold(n.Data, parts[0]) { return false }
    if len(parts) < 2 { return true }
    for _, a := range n.Attr {
        if a.Key == "class" {
            for _, class := range strings.Fields(a.Val) {
                if class == parts[1] { return true }
            }
        }
    }
    return false
}

func descendants(n *html.Node, selector string, output *[]*html.Node) {
    for child := n.FirstChild; child != nil; child = child.NextSibling {
        if matches(child, selector) { *output = append(*output, child) }
        descendants(child, selector, output)
    }
}

func children(n *html.Node, selector string, output *[]*html.Node) {
    for child := n.FirstChild; child != nil; child = child.NextSibling {
        if matches(child, selector) { *output = append(*output, child) }
    }
}

// Find selects matching descendants using the subset of selectors above.
func (s *Selection) Find(selector string) *Selection {
    if s == nil { return &Selection{} }
    tokens := strings.Fields(strings.ReplaceAll(selector, ">", " > "))
    current := s.nodes
    direct := false
    for _, token := range tokens {
        if token == ">" { direct = true; continue }
        next := make([]*html.Node, 0)
        seen := make(map[*html.Node]bool)
        for _, n := range current {
            var matches []*html.Node
            if direct { children(n, token, &matches) } else { descendants(n, token, &matches) }
            for _, m := range matches {
                if !seen[m] {
                    next = append(next, m)
                    seen[m] = true
                }
            }
        }
        current = next
        direct = false
    }
    return &Selection{nodes: current}
}

func (s *Selection) Each(fn func(int, *Selection)) *Selection {
    if s == nil { return s }
    for i, n := range s.nodes {
        fn(i, &Selection{nodes: []*html.Node{n}})
    }
    return s
}

func (s *Selection) First() *Selection {
    if s == nil || len(s.nodes) == 0 { return &Selection{} }
    return &Selection{nodes: s.nodes[:1]}
}

func (s *Selection) Attr(key string) (string, bool) {
    if s == nil || len(s.nodes) == 0 { return "", false }
    for _, a := range s.nodes[0].Attr {
        if a.Key == key { return a.Val, true }
    }
    return "", false
}

func appendText(b *strings.Builder, n *html.Node) {
    if n.Type == html.TextNode { b.WriteString(n.Data) }
    for child := n.FirstChild; child != nil; child = child.NextSibling {
        appendText(b, child)
    }
}
func (s *Selection) Text() string {
    if s == nil { return "" }
    var b strings.Builder
    for _, n := range s.nodes { appendText(&b, n) }
    return b.String()
}
