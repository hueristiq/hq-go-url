package url

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// Parser is responsible for parsing URLs with additional domain-specific information.
// It extends standard URL parsing to include extraction of subdomain, root domain, and TLD,
// while also providing functionality to add a default scheme if one is not provided.
// This is useful for handling URLs that are missing schemes and for accurately parsing the domain parts.
type Parser struct {
	scheme string // Default scheme to use if none is specified in the URL.

	dp *DomainParser // DomainParser is used for extracting domain-specific details.
}

// WithDefaultScheme sets the default scheme for the Parser, such as "http" or "https".
// This scheme will be applied to URLs that do not include a scheme when they are parsed.
func (up *Parser) WithDefaultScheme(scheme string) {
	ParserWithDefaultScheme(scheme)(up)
}

// DefaultScheme returns the currently set default scheme for the Parser.
// This is useful for checking what scheme will be added to URLs that do not specify one.
func (up *Parser) DefaultScheme() (scheme string) {
	return up.scheme
}

// Parse takes a raw URL string and parses it into a custom URL struct.
// This struct includes standard URL components (from net/url) as well as domain-specific details
// like subdomain, root domain, and TLD. If the URL does not include a scheme, the default scheme
// is added. The method also handles splitting the host and port if they are specified.
//
// Returns:
//   - parsedURL: A pointer to the parsed URL, containing all its components.
//   - err: An error if the URL is invalid or cannot be parsed.
func (up *Parser) Parse(rawURL string) (parsedURL *URL, err error) {
	parsedURL = &URL{}

	// Add default scheme if necessary
	if up.scheme != "" {
		rawURL = addScheme(rawURL, up.scheme)
	}

	// Standard URL parsing
	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("error parsing URL: %w", err)

		return
	}

	// Split host and port, and handle errors
	parsedURL.Host, parsedURL.Port, err = splitHostPort(parsedURL.Host)
	if err != nil {
		err = fmt.Errorf("error splitting host and port: %w", err)

		return
	}

	domainRegex := regexp.MustCompile(`(?i)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}`)

	if domainRegex.MatchString(parsedURL.Host) {
		parsedURL.Domain = up.dp.Parse(parsedURL.Host)
	}

	// Extract file extension from the path
	parsedURL.Extension = path.Ext(parsedURL.Path)

	return
}

// ParserOptionsFunc defines a function type for configuring a Parser instance.
// This allows setting options such as the default scheme or custom domain parsing logic.
type ParserOptionsFunc func(*Parser)

// ParserInterface defines the interface for URL parsing functionality.
// It ensures that any Parser implementation can set default schemes and parse URLs.
type ParserInterface interface {
	WithDefaultScheme(scheme string)                 // Set the default scheme.
	DefaultScheme() (scheme string)                  // Get the default scheme.
	Parse(rawURL string) (parsedURL *URL, err error) // Parse a raw URL string into a URL struct.
}

// Ensure that Parser implements the ParserInterface.
var _ ParserInterface = &Parser{}

// NewParser creates and initializes a new Parser with the given options.
// It also sets up a DomainParser for extracting domain-specific details (subdomain, root domain, and TLD).
// Additional configuration options can be applied using the variadic opts parameter.
//
// Returns:
//   - up: A pointer to the initialized Parser.
func NewParser(opts ...ParserOptionsFunc) (up *Parser) {
	up = &Parser{}

	// Initialize the DomainParser for domain-specific parsing.
	dp := NewDomainParser()
	up.dp = dp

	// Apply any additional options provided to configure the Parser.
	for _, opt := range opts {
		opt(up)
	}

	return
}

// ParserWithDefaultScheme returns a ParserOptionsFunc that sets the default scheme for the Parser.
// This function allows you to specify a scheme that will be used for URLs that don't provide one.
//
// Parameters:
//   - scheme: The default scheme to use (e.g., "http" or "https").
//
// Returns:
//   - A function that applies the default scheme to the Parser.
func ParserWithDefaultScheme(scheme string) ParserOptionsFunc {
	return func(up *Parser) {
		up.scheme = scheme
	}
}

// addScheme is a helper function that adds a scheme to a URL if it's missing.
// This function ensures that URLs without schemes are treated as absolute URLs
// rather than relative paths by adding the appropriate scheme prefix.
//
// Parameters:
//   - inURL: The raw URL string.
//   - scheme: The scheme to be added (if missing).
//
// Returns:
//   - outURL: The URL string with the scheme added (if necessary).
func addScheme(inURL, scheme string) (outURL string) {
	switch {
	case strings.HasPrefix(inURL, "//"): // If URL starts with "//", add the scheme.
		outURL = scheme + ":" + inURL
	case strings.HasPrefix(inURL, "://"): // If URL starts with "://", prepend the scheme.
		outURL = scheme + inURL
	case !strings.Contains(inURL, "//"): // If URL does not contain "//", prepend "://".
		outURL = scheme + "://" + inURL
	default:
		outURL = inURL
	}

	return
}

// splitHostPort separates the host and port components in a network address.
// This function handles both IPv4 and IPv6 addresses and gracefully manages URLs
// that do not include a port. Unlike net.SplitHostPort(), it retains the brackets
// around IPv6 addresses.
//
// Parameters:
//   - address: The raw address string containing the host and optional port.
//
// Returns:
//   - host: The host part of the address.
//   - port: The port part of the address, if present.
//   - err: An error if the port is not a valid integer.
func splitHostPort(address string) (host string, port int, err error) {
	host = address

	// Look for the last colon, which separates the host and port.
	i := strings.LastIndex(address, ":")
	if i == -1 {
		return
	}

	// Handle IPv6 addresses enclosed in brackets.
	if strings.HasPrefix(address, "[") && strings.Contains(address[i:], "]") {
		return
	}

	// Split the host and port.
	host = address[:i]

	// Convert the port part to an integer.
	if port, err = strconv.Atoi(address[i+1:]); err != nil {
		return
	}

	return
}
