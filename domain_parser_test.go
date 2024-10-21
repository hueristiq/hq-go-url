package url_test

import (
	"testing"

	hqgourl "github.com/hueristiq/hq-go-url"
	"github.com/stretchr/testify/assert"
)

// Test parsing of a valid domain with subdomain, SLD, and TLD.
func TestDomainParser_Parse_ValidDomain(t *testing.T) {
	t.Parallel()

	domain := "www.example.com"

	parser := hqgourl.NewDomainParser()

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "www", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "com", parsed.TLD)
}

// Test parsing of a domain without subdomain.
func TestDomainParser_Parse_DomainWithoutSubdomain(t *testing.T) {
	t.Parallel()

	domain := "example.com"

	parser := hqgourl.NewDomainParser()

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain) // No subdomain.
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "com", parsed.TLD)
}

// Test parsing of a domain without a valid TLD.
func TestDomainParser_Parse_InvalidTLD(t *testing.T) {
	t.Parallel()

	domain := "example.invalidtld"

	parser := hqgourl.NewDomainParser()

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)             // No subdomain.
	assert.Equal(t, "example.invalidtld", parsed.SLD) // Treat the whole domain as SLD.
	assert.Equal(t, "", parsed.TLD)
}

// Test parsing of a domain with a pseudo-TLD.
func TestDomainParser_Parse_PseudoTLD(t *testing.T) {
	t.Parallel()

	domain := "example.local"

	parser := hqgourl.NewDomainParser()

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "local", parsed.TLD) // Recognized pseudo-TLD.
}

// Test parsing of a single-word domain (no TLD or subdomain).
func TestDomainParser_Parse_SingleWordDomain(t *testing.T) {
	t.Parallel()

	domain := "localhost"

	parser := hqgourl.NewDomainParser()

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "localhost", parsed.SLD) // No TLD, treated as SLD.
	assert.Equal(t, "", parsed.TLD)
}

// Test parsing with custom TLDs.
func TestDomainParserWithCustomTLDs(t *testing.T) {
	t.Parallel()

	domain := "example.custom"

	parser := hqgourl.NewDomainParser(hqgourl.DomainParserWithTLDs("custom"))

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "custom", parsed.TLD) // Recognizes custom TLD.
}

// Test parsing an empty domain string.
func TestDomainParser_Parse_EmptyString(t *testing.T) {
	t.Parallel()

	domain := ""

	parser := hqgourl.NewDomainParser()

	parsed := parser.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "", parsed.SLD) // No SLD for an empty domain.
	assert.Equal(t, "", parsed.TLD)
}
