// Package extractor provides functionality for extracting URLs from text using highly configurable
// regular expressions. It is designed to support a wide range of URL formats including fully qualified
// web URLs, email addresses, and relative URLs. The package offers fine-grained control over the URL
// extraction process through the Extractor type, which can be customized using functional options to
// require specific URL components (such as schemes or hosts) or to supply custom regular expression
// patterns.
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
package extractor
