// Package extractor provides utilities for extracting URLs from text using configurable
// regular expressions. It allows users to define extraction rules, including constraints
// on URL schemes and hosts. The package supports standard and custom extraction patterns
// to accommodate various use cases, such as email extraction, relative URLs, and domain-based filtering.
//
// Features:
// - Configurable URL extraction with optional constraints (scheme, host, etc.).
// - Customizable regex patterns for URL matching.
// - Support for official, unofficial, and no-authority URL schemes.
// - Ability to handle ASCII and Unicode TLDs.
// - Extraction of IPv4, IPv6 addresses, and ports from URLs.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "go.source.hueristiq.com/url/extractor"
//	)
//
//	func main() {
//	    e := extractor.New(extractor.WithScheme())
//	    regex := e.CompileRegex()
//	    urls := regex.FindAllString("Visit https://example.com", -1)
//	    fmt.Println(urls) // Output: ["https://example.com"]
//	}
//
// This package is designed to be highly flexible, enabling developers to extract URLs
// from various sources while applying custom constraints as needed.
package extractor
