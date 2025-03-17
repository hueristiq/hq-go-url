# hq-go-url

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/hq-go-url)](https://goreportcard.com/report/github.com/hueristiq/hq-go-url) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-url.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-url/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-url.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-url/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-url/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-url/blob/master/CONTRIBUTING.md)

`hq-go-url` is a [Go (Golang)](http://golang.org/) package for extracting, parsing, and manipulating URLs with ease. It provides powerful tools to:

* Extract URLs from text using flexible, customizable regular expressions.
* Parse URLs with extended functionality that includes detailed domain analysis—splitting a domain into its subdomain, second-level domain (SLD), and top-level domain (TLD).

This package is especially useful for developers who need to work with URLs in a structured way.

## Resources

* [Features](#features)
* [Usage](#usage)
	* [Extraction](#extraction)
	* [Parsing](#parsing)
		* [Domains](#domains)
		* [URLs](#urls)
* [Contributing](#contributing)
* [Licensing](#licensing)

## Features

* **Configurable URL Extraction:** Extract URLs from text using regular expressions.
* **Extended URL Parsing:** Extend the standard [`net/url`](https://pkg.go.dev/net/url) package in Go with additional fields and capabilities.

## Usage

```bash
go get -v -u go.source.hueristiq.com/url
```

Below are examples demonstrating how to use the different features of the `hq-go-url` package.

### Extraction

The extractor package allows you to extract URLs from a block of text using regular expressions that can be customized to your needs.

```go
package main

import (
    "fmt"
    "log"

    "go.source.hueristiq.com/url/extractor"
)

func main() {
    // Create a new extractor with the requirement that URLs include a scheme.
    ex := extractor.New(extractor.WithScheme())

    // Compile the regular expression used for URL extraction.
    regex := ex.CompileRegex()

    text := "Visit our website at https://www.example.com or contact us at info@example.com."

    // Extract all matching URLs from the text.
    urls := regex.FindAllString(text, -1)

    fmt.Println("Extracted URLs:")

    for _, u := range urls {

        fmt.Println(u)
    }
}
```

You can customize how URLs are extracted by specifying URL schemes, hosts, or providing custom regular expression patterns.

* Extract URLs with Schemes Pattern:

	```go
	e := extractor.New(
		extractor.WithSchemePattern(`(?:https?|ftp)://`),
	)
	```

	This configuration will extract URLs with `http`, `https`, or `ftp` schemes.

* Extract URLs with Host Pattern:

	```go
	e := extractor.New(
		extractor.WithHostPattern(`(?:www\.)?example\.com`),
	)

	```

	This configuration will extract URLs that have hosts matching `www.example.com` or `example.com`.

### Parsing

```go
package main

import (
	"fmt"

	"go.source.hueristiq.com/url/parser"
)

func main() {
	p := parser.NewParser()

	parsed, err := p.Parse("https://subdomain.example.com:8080/path/file.txt")
	if err != nil {
		fmt.Println("Error parsing URL:", err)

		return
	}

	fmt.Printf("Scheme: %s\n", parsed.Scheme)
	fmt.Printf("Host: %s\n", parsed.Host)
	fmt.Printf("Hostname: %s\n", parsed.Hostname())
	fmt.Printf("Subdomain: %s\n", parsed.Domain.Subdomain)
	fmt.Printf("SLD: %s\n", parsed.Domain.SLD)
	fmt.Printf("TLD: %s\n", parsed.Domain.TLD)
	fmt.Printf("Port: %s\n", parsed.Port())
	fmt.Printf("Path: %s\n", parsed.Path)
}
```

You can customize how URLs are parsed by specifying default scheme, or providing custom TLDs.

* Parse URLs with default scheme:

	```go
	p := parser.NewParser(parser.WithDefaultScheme("https"))
	```

* Parse URLs with custom TLDs:

	```go
	p := parser.NewParser(parser.WithTLDs("custom", "custom2"))
	```

## Contributing

Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-url/pulls) or report [Issues](https://github.com/hueristiq/hq-go-url/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-url/blob/master/CONTRIBUTING.md).

Huge thanks to the [contributors](https://github.com/hueristiq/hq-go-url/graphs/contributors) thus far!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-url&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-url/blob/master/LICENSE).