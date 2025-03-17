// Package parser provides robust URL parsing with enhanced domain extraction capabilities.
// The parser package extends the standard net/url package by embedding its URL struct and
// augmenting it with detailed domain parsing. It decomposes a URL’s host into three key components:
//   - Subdomain (e.g., "www" in "www.example.com")
//   - Second-Level Domain (SLD, e.g., "example" in "www.example.com")
//   - Top-Level Domain (TLD, e.g., "com" in "www.example.com")
//
// Example Usage:
//
//		package main
//
//		import (
//		    "fmt"
//
//		    "go.source.hueristiq.com/url/parser"
//		)
//
//		func main() {
//		    parser := NewParser(WithDefaultScheme("https"))
//
//		    parsedURL, err := parser.Parse("www.example.com")
//		    if err != nil {
//		        // Handle error.
//		    }
//
//	     // Access the reconstructed domain.
//		    fmt.Println(parsedURL.Domain.String()) // Output: "www.example.com"
//		}
package parser
