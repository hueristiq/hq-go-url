package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.source.hueristiq.com/url/parser"
)

func Test_DomainParser_Parse(t *testing.T) {
	t.Parallel()

	p := parser.NewDomainParser()

	tests := []struct {
		name                 string
		domain               string
		expectedParsedDomain *parser.Domain
	}{
		{
			"domain",
			"example.com",
			&parser.Domain{
				Subdomain: "",
				SLD:       "example",
				TLD:       "com",
			},
		},
		{
			"domain with subdomain",
			"www.example.com",
			&parser.Domain{
				Subdomain: "www",
				SLD:       "example",
				TLD:       "com",
			},
		},
		{
			"domain with invalid TLD",
			"example.invalidtld",
			&parser.Domain{
				Subdomain: "",
				SLD:       "example.invalidtld",
				TLD:       "",
			},
		},
		{
			"domain with pseudo TLD",
			"example.local",
			&parser.Domain{
				Subdomain: "",
				SLD:       "example",
				TLD:       "local",
			},
		},
		{
			"domain (single word)",
			"localhost",
			&parser.Domain{
				Subdomain: "",
				SLD:       "localhost",
				TLD:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualParsedDomain := p.Parse(tt.domain)

			assert.Equal(t, tt.expectedParsedDomain, actualParsedDomain, "Expected and actual Person structs should be equal")
		})
	}
}

func Test_DomainParser_WithTLDs_Parse(t *testing.T) {
	t.Parallel()

	domain := "example.custom"

	p := parser.NewDomainParser(parser.DomainParserWithTLDs("custom"))

	parsed := p.Parse(domain)

	assert.NotNil(t, parsed)
	assert.Equal(t, "", parsed.Subdomain)
	assert.Equal(t, "example", parsed.SLD)
	assert.Equal(t, "custom", parsed.TLD)
}
