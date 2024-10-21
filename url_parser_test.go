package url_test

import (
	"testing"

	hqgourl "github.com/hueristiq/hq-go-url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test parsing a valid URL with a scheme and domain.
func TestParser_Parse_ValidURL(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser()

	parsed, err := parser.Parse("https://www.example.com/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify standard URL components.
	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "www.example.com", parsed.Host)
	assert.Equal(t, "/path", parsed.Path)

	// Verify domain components.
	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "www", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

// Test parsing an invalid URL.
func TestParser_Parse_InvalidURL(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser()

	_, err := parser.Parse("://example.com")

	require.Error(t, err)

	assert.Contains(t, err.Error(), "error parsing URL")
}

// Test parsing a URL without a scheme and adding the default scheme.
func TestParser_Parse_URLWithoutScheme(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser(hqgourl.ParserWithDefaultScheme("https"))

	parsed, err := parser.Parse("example.com/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify that the default scheme has been added.
	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/path", parsed.Path)

	// Verify domain components.
	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain) // No subdomain
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

// Test parsing a URL with a subdomain.
func TestParser_Parse_URLWithSubdomain(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser()

	parsed, err := parser.Parse("https://sub.example.com/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify domain components.
	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "sub", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

// Test parsing a URL with a port number.
func TestParser_Parse_URLWithPort(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser()

	parsed, err := parser.Parse("https://example.com:8080/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify URL components.
	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "example.com:8080", parsed.Host)
	assert.Equal(t, "example.com", parsed.Hostname()) // Hostname should exclude the port.
	assert.Equal(t, "/path", parsed.Path)

	// Verify domain components.
	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

// Test parsing a URL with a custom scheme.
func TestParser_Parse_CustomScheme(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser(hqgourl.ParserWithDefaultScheme("ftp"))

	parsed, err := parser.Parse("example.com/file.txt")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify that the custom default scheme has been added.
	assert.Equal(t, "ftp", parsed.Scheme)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/file.txt", parsed.Path)

	// Verify domain components.
	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

// Test parsing a URL with a scheme already specified.
func TestParser_Parse_AlreadyHasScheme(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser(hqgourl.ParserWithDefaultScheme("https"))

	parsed, err := parser.Parse("http://example.com/file.txt")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Ensure that the existing scheme (http) is not replaced by the default scheme (https).
	assert.Equal(t, "http", parsed.Scheme)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/file.txt", parsed.Path)

	// Verify domain components.
	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

// Test parsing a URL with an IPv4 address.
func TestParser_Parse_URLWithIPv4Address(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser()

	parsed, err := parser.Parse("http://192.168.0.1/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify standard URL components.
	assert.Equal(t, "http", parsed.Scheme)
	assert.Equal(t, "192.168.0.1", parsed.Host)
	assert.Equal(t, "/path", parsed.Path)

	// Ensure that the domain parsing doesn't apply to IP addresses.
	assert.Nil(t, parsed.Domain)
}

// Test parsing a URL with an IPv6 address.
func TestParser_Parse_URLWithIPv6Address(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewParser()

	parsed, err := parser.Parse("https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	// Verify standard URL components.
	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000", parsed.Host)
	assert.Equal(t, "", parsed.Path)

	// Ensure that the domain parsing doesn't apply to IP addresses.
	assert.Nil(t, parsed.Domain)
}
