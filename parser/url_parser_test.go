package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.source.hueristiq.com/url/parser"
)

func TestParser_Parse_ValidURL(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser()

	parsed, err := p.Parse("https://www.example.com/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "www.example.com", parsed.Host)
	assert.Equal(t, "/path", parsed.Path)

	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "www", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

func TestParser_Parse_InvalidURL(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser()

	_, err := p.Parse("://example.com")

	require.Error(t, err)

	assert.Contains(t, err.Error(), "error parsing URL")
}

func TestParser_Parse_URLWithoutScheme(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser(parser.URLParserWithDefaultScheme("https"))

	parsed, err := p.Parse("example.com/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/path", parsed.Path)

	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

func TestParser_Parse_URLWithSubdomain(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser()

	parsed, err := p.Parse("https://sub.example.com/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "sub", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

func TestParser_Parse_URLWithPort(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser()

	parsed, err := p.Parse("https://example.com:8080/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "example.com:8080", parsed.Host)
	assert.Equal(t, "example.com", parsed.Hostname())
	assert.Equal(t, "/path", parsed.Path)

	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

func TestParser_Parse_CustomScheme(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser(parser.URLParserWithDefaultScheme("ftp"))

	parsed, err := p.Parse("example.com/file.txt")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "ftp", parsed.Scheme)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/file.txt", parsed.Path)

	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

func TestParser_Parse_AlreadyHasScheme(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser(parser.URLParserWithDefaultScheme("https"))

	parsed, err := p.Parse("http://example.com/file.txt")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "http", parsed.Scheme)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/file.txt", parsed.Path)

	assert.NotNil(t, parsed.Domain)
	assert.Equal(t, "", parsed.Domain.Subdomain)
	assert.Equal(t, "example", parsed.Domain.SLD)
	assert.Equal(t, "com", parsed.Domain.TLD)
}

func TestParser_Parse_URLWithIPv4Address(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser()

	parsed, err := p.Parse("http://192.168.0.1/path")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "http", parsed.Scheme)
	assert.Equal(t, "192.168.0.1", parsed.Host)
	assert.Equal(t, "/path", parsed.Path)

	assert.Nil(t, parsed.Domain)
}

func TestParser_Parse_URLWithIPv6Address(t *testing.T) {
	t.Parallel()

	p := parser.NewURLParser()

	parsed, err := p.Parse("https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000")

	require.NoError(t, err)

	assert.NotNil(t, parsed)

	assert.Equal(t, "https", parsed.Scheme)
	assert.Equal(t, "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000", parsed.Host)
	assert.Equal(t, "", parsed.Path)

	assert.Nil(t, parsed.Domain)
}
