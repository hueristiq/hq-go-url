package parser_test

import (
	"net/url"
	"testing"

	"github.com/hueristiq/hq-go-url/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parser_Parse(t *testing.T) {
	t.Parallel()

	p := parser.New()

	tests := []struct {
		name           string
		raw            string
		expectedParsed *parser.URL
		expectedErr    bool
	}{
		{
			"URL",
			"https://example.com/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
		{
			"URL with subdomain",
			"https://www.example.com/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "www.example.com",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "www",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
		{
			"URL with invalid TLD",
			"https://www.example.invalidtld/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "www.example.invalidtld",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "www.example.invalidtld",
					TopLevelDomain:    "",
				},
			},
			false,
		},
		{
			"URL with pseudo  TLD",
			"https://www.example.local/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "www.example.local",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "www",
					SecondLevelDomain: "example",
					TopLevelDomain:    "local",
				},
			},
			false,
		},
		{
			"URL with port",
			"https://www.example.com:8080/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "www.example.com:8080",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "www",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
		{
			"URL with IPv4",
			"http://192.168.0.1/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "192.168.0.1",
					Path:   "/path",
				},
				Domain: nil,
			},
			false,
		},
		{
			"URL with IPv6",
			"https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000",
					Path:   "/path",
				},
				Domain: nil,
			},
			false,
		},
		{
			"Invalid URL 0",
			"://example.com",
			&parser.URL{},
			true,
		},
		{
			"Invalid URL 1",
			"https://example.com/%invalid",
			&parser.URL{},
			true,
		},
		{
			"Invalid URL 2 <- Path get's normalized",
			"https://example.com/w%0d%2e/",
			&parser.URL{
				URL: &url.URL{
					Scheme:  "https",
					Host:    "example.com",
					Path:    "/w\r./",
					RawPath: "/w%0d%2e/",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
		{
			"Invalid URL 3",
			"example.com/with/path?some'param=`'+OR+ORDER+BY+1--",
			&parser.URL{
				URL: &url.URL{
					Scheme:   "",
					Host:     "",
					Path:     "example.com/with/path",
					RawQuery: "some'param=`'+OR+ORDER+BY+1--",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "",
					TopLevelDomain:    "",
				},
			},
			false,
		},
		{
			"Invalid URL 4 <- Path get's normalized",
			"https://example.com/a/b/c/../c",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/a/b/c/../c",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualParsedURL, err := p.Parse(tt.raw)

			if tt.expectedErr {
				require.Error(t, err, "Expected an error but got none")

				return
			}

			assert.Equal(t, tt.expectedParsed, actualParsedURL, "Expected and actual Person structs should be equal")
		})
	}
}

func Test_Parser_WithDefaultScheme_Parse(t *testing.T) {
	t.Parallel()

	p := parser.New(parser.WithDefaultScheme("https"))

	tests := []struct {
		name              string
		URL               string
		expectedParsedURL *parser.URL
		expectedErr       bool
	}{
		{
			"URL without scheme",
			"example.com/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
		{
			"URL with ://",
			"://example.com/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
		{
			"URL with scheme",
			"http://example.com/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "com",
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualParsedURL, err := p.Parse(tt.URL)

			if tt.expectedErr {
				require.Error(t, err, "Expected an error but got none")

				return
			}

			assert.Equal(t, tt.expectedParsedURL, actualParsedURL, "Expected and actual Person structs should be equal")
		})
	}
}

func Test_Parser_WithTLDs_Parse(t *testing.T) {
	t.Parallel()

	p := parser.New(parser.WithTLDs("custom"))

	tests := []struct {
		name              string
		URL               string
		expectedParsedURL *parser.URL
		expectedErr       bool
	}{
		{
			"URL without scheme",
			"https://example.custom/path",
			&parser.URL{
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.custom",
					Path:   "/path",
				},
				Domain: &parser.Domain{
					Subdomain:         "",
					SecondLevelDomain: "example",
					TopLevelDomain:    "custom",
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualParsedURL, err := p.Parse(tt.URL)

			if tt.expectedErr {
				require.Error(t, err, "Expected an error but got none")

				return
			}

			assert.Equal(t, tt.expectedParsedURL, actualParsedURL, "Expected and actual Person structs should be equal")
		})
	}
}
