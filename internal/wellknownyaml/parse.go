// Package wellknownyaml parses the fixed-format YAML data file used by the
// well-known RPC endpoint generator. It is not a general-purpose YAML parser.
package wellknownyaml

import (
    "bufio"
    "bytes"
    "fmt"
    "strconv"
    "strings"
)

type Entry struct {
    UUID string
    Domain string
    Annotation string
    IfName string
    IfVersions []string
    Endpoints []string
    Keywords []string
}

func scalar(s string) (string, error) {
    s = strings.TrimSpace(s)
    if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
        return strings.ReplaceAll(s[1:len(s)-1], "''", "'"), nil
    }
    if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
        return strconv.Unquote(s)
    }
    if strings.HasPrefix(s, "|") || strings.HasPrefix(s, ">") {
        return "", fmt.Errorf("block YAML strings are not supported in endpoint data")
    }
    return s, nil
}

func list(s string) ([]string, error) {
    s = strings.TrimSpace(s)
    if s == "[]" { return nil, nil }
    if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
        return nil, fmt.Errorf("invalid YAML flow sequence %q", s)
    }
    s = s[1:len(s)-1]
    var out []string
    var b strings.Builder
    var quote byte
    for i := 0; i < len(s); i++ {
        c := s[i]
        if quote != 0 {
            b.WriteByte(c)
            if c == quote {
                if quote == '\'' && i+1 < len(s) && s[i+1] == '\'' {
                    i++
                    b.WriteByte(s[i])
                } else {
                    quote = 0
                }
            } else if quote == '"' && c == '\\' && i+1 < len(s) {
                i++
                b.WriteByte(s[i])
            }
            continue
        }
        if c == '\'' || c == '"' { quote = c; b.WriteByte(c); continue }
        if c != ',' { b.WriteByte(c); continue }
        item, err := scalar(b.String())
        if err != nil { return nil, err }
        out = append(out, item)
        b.Reset()
    }
    if quote != 0 { return nil, fmt.Errorf("unterminated flow-sequence quote") }
    if strings.TrimSpace(b.String()) != "" {
        item, err := scalar(b.String())
        if err != nil { return nil, err }
        out = append(out, item)
    }
    return out, nil
}

func Unmarshal(data []byte, dst *[]Entry) error {
    if dst == nil { return fmt.Errorf("nil endpoint destination") }
    scanner := bufio.NewScanner(bytes.NewReader(data))
    scanner.Buffer(make([]byte, 4096), 4*1024*1024)
    result := make([]Entry, 0)
    active := ""
    field := func(e *Entry, key, value string) error {
        value = strings.TrimSpace(value)
        switch key {
        case "uuid":
            v, err := scalar(value); if err != nil { return err }; e.UUID = v
        case "domain":
            v, err := scalar(value); if err != nil { return err }; e.Domain = v
        case "annotation":
            v, err := scalar(value); if err != nil { return err }; e.Annotation = v
        case "if_name":
            v, err := scalar(value); if err != nil { return err }; e.IfName = v
        case "if_versions", "endpoints", "keywords":
            if value == "" { active = key; return nil }
            items, err := list(value); if err != nil { return err }
            switch key {
            case "if_versions": e.IfVersions = items
            case "endpoints": e.Endpoints = items
            case "keywords": e.Keywords = items
            }
        }
        return nil
    }
    for scanner.Scan() {
        line := strings.TrimRight(scanner.Text(), "\r")
        if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") { continue }
        if strings.HasPrefix(line, "- uuid:") {
            result = append(result, Entry{})
            active = ""
            if err := field(&result[len(result)-1], "uuid", strings.TrimPrefix(line, "- uuid:")); err != nil { return err }
            continue
        }
        if len(result) == 0 { return fmt.Errorf("expected sequence entry before %q", line) }
        if strings.HasPrefix(line, "  - ") {
            if active == "" { return fmt.Errorf("unexpected YAML sequence item %q", line) }
            item, err := scalar(strings.TrimPrefix(line, "  - "))
            if err != nil { return err }
            e := &result[len(result)-1]
            switch active {
            case "endpoints": e.Endpoints = append(e.Endpoints, item)
            case "if_versions": e.IfVersions = append(e.IfVersions, item)
            case "keywords": e.Keywords = append(e.Keywords, item)
            }
            continue
        }
        if !strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "    ") {
            return fmt.Errorf("unsupported endpoint YAML indentation: %q", line)
        }
        key, value, ok := strings.Cut(strings.TrimSpace(line), ":")
        if !ok { return fmt.Errorf("invalid endpoint YAML key: %q", line) }
        active = ""
        if err := field(&result[len(result)-1], key, value); err != nil { return err }
    }
    if err := scanner.Err(); err != nil { return err }
    for _, e := range result {
        if e.UUID == "" { return fmt.Errorf("entry missing UUID") }
    }
    *dst = result
    return nil
}
