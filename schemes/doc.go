// Package schemes provides categorized lists of URL schemes, including official IANA-assigned schemes,
// commonly used unofficial schemes, and schemes that do not require an authority component.
//
// URL schemes define how a resource should be accessed. They serve as the prefix in URLs, determining the
// protocol to be used when retrieving or interacting with a resource. This package helps with processing,
// validating, and handling different types of schemes in web applications, networking tools, and other
// software components.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/hueristiq/hq-go-url/schemes"
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
