// Package parser provides comprehensive functionality for parsing both domain names and URLs.
// It includes two primary sets of tools:
//
//  1. Domain Parsing:
//     The package defines types and methods to decompose a full domain name (e.g., "www.example.com")
//     into its constituent parts:
//     - Subdomain: e.g., "www"
//     - SLD (Second-Level Domain): e.g., "example"
//     - TLD (Top-Level Domain): e.g., "com"
//     This functionality is encapsulated in the Domain struct and the DomainParser type. The DomainParser
//     leverages a suffix array built from a list of known TLDs (both official and pseudo-TLDs) to quickly
//     and accurately determine where the TLD begins in a domain name. This allows for effective domain
//     validation, classification, and further analysis.
//
//  2. URL Parsing:
//     The package extends standard URL parsing (from net/url) by defining a custom URL type that embeds
//     *url.URL and adds a Domain field. The URLParser type is responsible for parsing raw URL strings into
//     this extended URL structure. In addition to the usual URL components (scheme, host, path, etc.),
//     URLParser also extracts detailed domain information using a DomainParser. It further provides the
//     capability to add a default scheme to URLs that are missing one, ensuring that relative URLs are
//     correctly interpreted as absolute.
//
// Features:
//
//   - **Domain Extraction:**
//     Split a domain string into subdomain, SLD, and TLD using efficient suffix array lookups.
//
//   - **Enhanced URL Parsing:**
//     Parse URLs with the added benefit of detailed domain breakdown, and automatically add default schemes
//     when needed.
//
//   - **Customization Options:**
//     Both the DomainParser and URLParser support configuration via functional options (e.g., setting
//     custom TLD lists or default schemes).
//
// Usage Example:
//
//	// Create a new DomainParser (using default TLDs).
//	domainParser := NewDomainParser()
//	// Parse a domain into its components.
//	parsedDomain := domainParser.Parse("www.example.com")
//	fmt.Println("Subdomain:", parsedDomain.Subdomain) // Output: "www"
//	fmt.Println("SLD:", parsedDomain.SLD)             // Output: "example"
//	fmt.Println("TLD:", parsedDomain.TLD)             // Output: "com"
//
//	// Create a URLParser with a default scheme of "https".
//	urlParser := NewURLParser(URLParserWithDefaultScheme("https"))
//	// Parse a raw URL. If the URL is missing a scheme, "https" will be prepended.
//	parsedURL, err := urlParser.Parse("www.example.com/path")
//	if err != nil {
//	    log.Fatalf("Error parsing URL: %v", err)
//	}
//	fmt.Println("Parsed URL:", parsedURL.String())
//	if parsedURL.Domain != nil {
//	    fmt.Println("Domain:", parsedURL.Domain.String()) // Output: "www.example.com"
//	}
//
// The parser package is a versatile and comprehensive toolset for URL and domain parsing. By
// combining robust domain extraction with extended URL parsing capabilities, it provides developers
// with the tools necessary to perform detailed analysis, validation, and manipulation of web addresses.
// This package is especially useful in contexts such as SEO analysis, security filtering, and any
// application where understanding the structure of a URL is important.
package parser
