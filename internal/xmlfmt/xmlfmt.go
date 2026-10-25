// Package xmlfmt formats XML using encoding/xml instead of a third-party library.
package xmlfmt

import (
    "bytes"
    "encoding/xml"
    "io"
    "strings"
)

// FormatXML returns an indented XML representation. If the input is not
// well-formed XML, it is returned unchanged rather than silently discarding it.
func FormatXML(input, prefix, indent string) string {
    decoder := xml.NewDecoder(strings.NewReader(input))
    decoder.Strict = false
    var output bytes.Buffer
    encoder := xml.NewEncoder(&output)
    encoder.Indent(prefix, indent)
    for {
        token, err := decoder.Token()
        if err == io.EOF { break }
        if err != nil { return input }
        if err := encoder.EncodeToken(token); err != nil { return input }
    }
    if err := encoder.Flush(); err != nil { return input }
    return output.String()
}
