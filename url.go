package url

import "net/url"

// URL extends the standard net/url URL struct by embedding it and adding additional fields
// for handling domain-related information. This extension provides a more detailed representation
// of the URL by including a separate `Domain` struct that breaks down the domain into Subdomain,
// second-level domain (SLD), and top-level domain (TLD).
//
// Fields:
//
//   - URL (*url.URL):
//
//   - Embeds the standard `net/url.URL` struct, which provides all the base URL parsing and
//     functionalities, such as handling the scheme, host, path, query parameters, and fragment.
//
//   - Methods and functions from the embedded `net/url.URL` can be used transparently.
//
//   - Domain (*Domain):
//
//   - A pointer to the `Domain` struct that contains parsed domain information, including:
//
//   - Subdomain (string): The subdomain of the URL (e.g., "www" in "www.example.com").
//
//   - Second-level domain (SLD) (string): The main domain (e.g., "example").
//
//   - Top-level domain (TLD) (string): The domain suffix (e.g., "com" in "www.example.com").
//
//   - This allows for better handling of domain components, which is useful in cases like:
//
//   - URL classification and domain analysis.
//
//   - Security or SEO applications where separating domain components is important.
//
// Example Usage:
//
//	// Parse a URL using the standard url.Parse method.
//	parsedURL, _ := url.Parse("https://www.example.com")
//
//	// Create an extended URL object and manually add domain information.
//	extendedURL := &URL{
//	    URL: parsedURL, // Embeds the parsed URL from the standard library.
//
//	    // Domain can be parsed separately or manually assigned.
//	    Domain: &Domain{
//	        Subdomain:      "www",     // Subdomain part (e.g., "www").
//	        SLD:     "example", // Root domain part (e.g., "example").
//	        TLD: "com",     // Top-level domain part (e.g., "com").
//	    },
//	}
//
//	// Access standard URL components.
//	fmt.Println(extendedURL.Scheme)   // Output: https
//	fmt.Println(extendedURL.Host)     // Output: www.example.com
//	fmt.Println(extendedURL.Path)     // Output: /
//
//	// Access domain-specific information.
//	fmt.Println(extendedURL.Domain.Subdomain)      // Output: www
//	fmt.Println(extendedURL.Domain.SLD)     // Output: example
//	fmt.Println(extendedURL.Domain.TLD) // Output: com
//
// Purpose:
//
//	This `URL` struct provides a more detailed breakdown of a URL's domain components,
//	making it particularly useful for tasks involving domain analysis, URL classification,
//	or scenarios where understanding subdomains, root domains, and TLDs is important.
type URL struct {
	*url.URL

	Domain *Domain
}
