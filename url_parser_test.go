package url_test

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	hqgourl "github.com/hueristiq/hq-go-url"
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	up := hqgourl.NewParser()

	if up == nil {
		t.Error("NewParser() = nil; want non-nil")
	}

	scheme := "https"

	parserWithDefaultScheme := hqgourl.NewParser(hqgourl.ParserWithDefaultScheme(scheme))

	if up == nil {
		t.Errorf("NewParser(ParserWithDefaultScheme(%s)) = nil; want non-nil", scheme)
	}

	expectedDefaultScheme := parserWithDefaultScheme.DefaultScheme()

	if expectedDefaultScheme != scheme {
		t.Errorf("NewParser(ParserWithDefaultScheme(%s)).DefaultScheme() = '%s', want '%s'", scheme, expectedDefaultScheme, scheme)
	}
}

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rawURL            string
		defaultScheme     string
		expectedParsedURL *hqgourl.URL
		expectParseErr    bool
	}{
		{
			"http://example.com",
			"http",
			&hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				Domain: &hqgourl.Domain{
					Sub:      "",
					Root:     "example",
					TopLevel: "com",
				},
			},
			false,
		},
		{
			"example.com",
			"http",
			&hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				Domain: &hqgourl.Domain{
					Sub:      "",
					Root:     "example",
					TopLevel: "com",
				},
			},
			false,
		},
		{
			"http://example.com/path/file.html",
			"http",
			&hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/path/file.html",
				},
				Domain: &hqgourl.Domain{
					Sub:      "",
					Root:     "example",
					TopLevel: "com",
				},
				Extension: ".html",
			},
			false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Parse(%q)", c.rawURL), func(t *testing.T) {
			t.Parallel()

			up := hqgourl.NewParser(
				hqgourl.ParserWithDefaultScheme(c.defaultScheme),
			)

			parsedURL, err := up.Parse(c.rawURL)

			if (err != nil) != c.expectParseErr {
				t.Errorf("Parse(%q) error = %v, expectParseErr %v", c.rawURL, err, c.expectParseErr)

				return
			}

			if !reflect.DeepEqual(parsedURL, c.expectedParsedURL) {
				t.Errorf("Parse(%q) = %+v, want %+v", c.rawURL, parsedURL, c.expectedParsedURL)
			}
		})
	}
}
