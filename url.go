package url

import "net/url"

//go:generate go run gen/schemes/main.go -output ./schemes/schemes_official.go
//go:generate go run gen/TLDs/main.go -output ./tlds/tlds_official.go
//go:generate go run gen/unicodes/main.go -output ./unicodes/unicodes.go

// URL extends the standard net/url URL struct by embedding it and adding additional domain-related fields.
// The extended struct includes subdomain, root domain, and top-level domain (TLD) information,
// along with standard URL components (such as scheme, host, and path) provided by the embedded URL struct.
//
// Additional Fields:
//   - Domain: A pointer to a Domain struct that contains parsed subdomain, root domain, and TLD details.
//   - Port: The port number used in the URL, if specified.
//   - Extension: The file extension derived from the URL's path, useful for identifying file types in URLs.
//
// This extended struct provides a more comprehensive representation of a URL, especially in scenarios
// where parsing and understanding domain parts (subdomains, root domains, and TLDs) is important.
//
// Example Usage:
//
//	parsedURL, _ := url.Parse("https://www.example.com:8080/index.html")
//	extendedURL := &URL{
//	    URL: parsedURL,
//	    Port: 8080,
//	    Extension: "html",
//	    Domain: &Domain{
//	        Sub: "www",
//	        Root: "example",
//	        TopLevel: "com",
//	    },
//	}
type URL struct {
	*url.URL // Embedding the standard net/url URL struct for base functionalities.

	Domain    *Domain // Contains parsed domain parts: subdomain, root domain, and TLD.
	Port      int     // Port number used in the URL.
	Extension string  // File extension derived from the URL path.
}
