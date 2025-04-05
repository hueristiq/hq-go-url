// Package tlds provides categorized lists of top-level domains (TLDs) and effective TLDs (eTLDs),
// including official IANA-approved domains and commonly used pseudo-TLDs.
//
// TLDs represent the highest level in the hierarchical domain name system, while eTLDs include
// public suffixes such as country code second-level domains (e.g., "co.uk", "gov.in") used for
// domain registrations. Pseudo-TLDs are unofficial but widely recognized domain suffixes used in
// specialized networks or alternative domain name systems.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/hueristiq/hq-go-url/tlds"
//	)
//
//	func main() {
//	    fmt.Println("Official TLDs:", tlds.Official)
//	    fmt.Println("Pseudo TLDs:", tlds.Pseudo)
//	}
//
// References:
// - IANA TLDs: https://data.iana.org/TLD/tlds-alpha-by-domain.txt
// - Public Suffix List: https://publicsuffix.org/list/public_suffix_list.dat
// - A-Z TLDs: https://tld-list.com/tlds-from-a-z
// - Wikipedia: https://en.wikipedia.org/wiki/Top-level_domain
// - Pseudo-TLDs: https://en.wikipedia.org/wiki/Pseudo-top-level_domain
package tlds
