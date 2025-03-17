// Package parser provides advanced URL parsing capabilities by extending the standard net/url package
// with additional domain extraction functionality.
//
// This package decomposes a URL’s host component into three primary parts:
//   - Subdomain (e.g., "www" in "www.example.com")
//   - Second-Level Domain (SLD, e.g., "example" in "www.example.com")
//   - Top-Level Domain (TLD, e.g., "com" in "www.example.com")
//
// The custom URL type embeds net/url.URL so that it integrates seamlessly with existing HTTP libraries,
// while the additional Domain struct holds the parsed components. The Parser type offers methods to parse
// raw URL strings into this extended URL struct. It also supports applying a default scheme when missing,
// and uses a suffix array for efficient TLD lookups.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "go.source.hueristiq.com/url/parser"
//	)
//
//	func main() {
//	    // Create a new parser with a default scheme of "https".
//	    p := parser.NewParser(parser.WithDefaultScheme("https"))
//
//	    // Parse a raw URL string without a scheme.
//	    parsedURL, err := p.Parse("www.example.com")
//	    if err != nil {
//	        fmt.Println("Error parsing URL:", err)
//	        return
//	    }
//
//	    // Print the reconstructed domain.
//	    fmt.Println("Domain:", parsedURL.Domain.String())
//	}
//
// References:
// - net/url package documentation: https://pkg.go.dev/net/url
package parser
