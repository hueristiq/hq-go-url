// Package schemes provides categorized lists of URL schemes, including official IANA-assigned schemes,
// commonly used unofficial schemes, and schemes that do not require an authority component.
//
// URL schemes define how a resource should be accessed. They serve as the prefix in URLs, determining the
// protocol to be used when retrieving or interacting with a resource. This package helps with processing,
// validating, and handling different types of schemes in web applications, networking tools, and other
// software components.
//
// Categories:
// - Official: Includes IANA-registered schemes, updated periodically from the official registry.
// - Unofficial: Consists of widely used but unregistered schemes specific to certain applications or services.
// - NoAuthority: Contains schemes that do not require an authority component and use a colon (":") instead of "://".
//
// Features:
// - Maintains structured lists of URL schemes for validation and processing.
// - Useful in security applications, protocol handlers, and URL parsing utilities.
// - References official IANA data sources to ensure accuracy.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "go.source.hueristiq.com/url/schemes"
//	)
//
//	func main() {
//	    fmt.Println("Supported schemes:", schemes.Official)
//	}
//
// References:
// - IANA Registry: https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml
// - Wikipedia URI Schemes: https://en.wikipedia.org/wiki/List_of_URI_schemes
package schemes
