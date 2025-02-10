// Package extractor provides functionality for extracting URLs from text using highly configurable
// regular expressions. It is designed to support a wide range of URL formats including fully qualified
// web URLs, email addresses, and relative URLs. The package offers fine-grained control over the URL
// extraction process through the Extractor type, which can be customized using functional options to
// require specific URL components (such as schemes or hosts) or to supply custom regular expression
// patterns.
//
// # Overview
//
// The primary type in this package is Extractor. An Extractor instance encapsulates configuration
// options for URL extraction. It allows you to control:
//
//   - Whether a URL scheme (e.g., "http", "https", "mailto") is mandatory.
//   - Whether a URL host (e.g., domain names or IP addresses) is mandatory.
//   - Custom regular expression patterns for matching URL schemes and hosts.
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
// The extractor package is a powerful tool for developers who need to accurately extract URLs from text.
// Its highly configurable nature makes it adaptable to a variety of application requirements, whether
// you are processing web content, parsing emails, or handling relative links. By combining default
// patterns with user-specified customizations, Extractor offers both robustness and flexibility in URL extraction.
package extractor
