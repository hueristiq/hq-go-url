package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.source.hueristiq.com/url/parser"
)

func TestDomainParser_Parse_ValidDomain(t *testing.T) {
	t.Parallel()

	domain := "www.example.com"

	p := parser.NewDomainParser()

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "www", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "com", parsed.TLD)
}

func TestDomainParser_Parse_DomainWithoutSubdomain(t *testing.T) {
	t.Parallel()

	domain := "example.com"

	p := parser.NewDomainParser()

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "com", parsed.TLD)
}

func TestDomainParser_Parse_InvalidTLD(t *testing.T) {
	t.Parallel()

	domain := "example.invalidtld"

	p := parser.NewDomainParser()

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example.invalidtld", parsed.SLD)
	assert.Equal(t, "", parsed.TLD)
}

func TestDomainParser_Parse_PseudoTLD(t *testing.T) {
	t.Parallel()

	domain := "example.local"

	p := parser.NewDomainParser()

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "local", parsed.TLD)
}

func TestDomainParser_Parse_SingleWordDomain(t *testing.T) {
	t.Parallel()

	domain := "localhost"

	p := parser.NewDomainParser()

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "localhost", parsed.SLD)
	assert.Equal(t, "", parsed.TLD)
}

func TestDomainParserWithCustomTLDs(t *testing.T) {
	t.Parallel()

	domain := "example.custom"

	p := parser.NewDomainParser(parser.DomainParserWithTLDs("custom"))

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "custom", parsed.TLD)
}

func TestDomainParser_Parse_EmptyString(t *testing.T) {
	t.Parallel()

	domain := ""

	p := parser.NewDomainParser()

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "", parsed.SLD)
	assert.Equal(t, "", parsed.TLD)
}
