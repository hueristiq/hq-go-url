// Package parser provides utilities for parsing and analyzing URLs, including detailed domain extraction.
//
// This package extends the standard net/url parsing capabilities by incorporating domain-specific parsing.
// It introduces the `URL` struct, which embeds `net/url.URL` and adds a `Domain` field containing subdomain,
// second-level domain (SLD), and top-level domain (TLD) information.
//
// Features:
//   - URL Parsing: Converts raw URL strings into structured `URL` objects.
//   - Domain Extraction: Breaks down domain names into subdomains, root domains, and TLDs.
//   - Customizable Parsing: Allows configuration via `URLParserOption` functions.
//
// Example Usage:
//
//	parser := NewURLParser(URLParserWithDefaultScheme("https"))
//	parsedURL, err := parser.Parse("example.com/path")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(parsedURL.Scheme)   // Output: https
//	fmt.Println(parsedURL.Hostname()) // Output: example.com
//	fmt.Println(parsedURL.Domain.SLD) // Output: example
package parser
