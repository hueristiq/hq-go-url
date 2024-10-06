package url

import "net/url"

//go:generate go run gen/schemes/main.go -output ./schemes/schemes_official.go
//go:generate go run gen/TLDs/main.go -output ./tlds/tlds_official.go
//go:generate go run gen/unicodes/main.go -output ./unicodes/unicodes.go

// URL extends the standard net/url URL struct with additional domain-related fields.
// It includes details like subdomain, root domain, and Top-Level Domain (TLD), along with
// standard URL components. This struct provides a comprehensive representation of a URL.
type URL struct {
	*url.URL // Embedding the standard URL struct for base functionalities.

	Domain    *Domain
	Port      int    // Port number used in the URL.
	Extension string // File extension derived from the URL path.
}
