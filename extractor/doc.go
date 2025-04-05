// Package extractor provides advanced URL extraction capabilities from text.
//
// URL extraction is a common requirement in text processing, data mining, and content analysis.
// This package offers a highly configurable extractor that uses a composite regular expression to
// identify and capture URLs in various forms. It supports fully-qualified URLs (with schemes and hosts),
// email addresses, and relative URLs, while allowing users to enforce or relax requirements for URL schemes
// and hosts. Custom regular expression patterns can also be provided to further fine-tune the extraction process.
//
// The extractor leverages robust Unicode and punycode handling, and it incorporates known TLD lists and
// scheme definitions (both official and unofficial) to ensure accurate matching of web addresses and email formats.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/hueristiq/hq-go-url/extractor"
//	)
//
//	func main() {
//	    // Create a new extractor that requires URL schemes.
//	    ext := extractor.New(extractor.WithScheme())
//
//	    // Compile the regex pattern based on the extractor configuration.
//	    regex := ext.CompileRegex()
//
//	    text := "Contact us at info@example.com or visit https://www.example.com for more details."
//	    urls := regex.FindAllString(text, -1)
//	    fmt.Println("Extracted URLs:", urls)
//	}
//
// References:
// - Regular Expression HOWTO: https://golang.org/pkg/regexp/
// - IANA URI Schemes: https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml
// - Unicode and UTF-8 handling in Go: https://golang.org/pkg/unicode/utf8/
package extractor
