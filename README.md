# hq-go-url

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/hq-go-url)](https://goreportcard.com/report/github.com/hueristiq/hq-go-url) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-url.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-url/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-url.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-url/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-url/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-url/blob/master/CONTRIBUTING.md)

`hq-go-url` is a [Go (Golang)](http://golang.org/) package for working with URLs. It provides robust tools for extracting URLs from text and parsing them into granular components.

## Resources

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
	- [Extraction](#extraction)
	- [Parsing](#parsing)
- [Contributing](#contributing)
- [Licensing](#licensing)

## Features

- **URL Extraction:** Utilizes advanced regular expression patterns to scan and extract valid URLs from any text. 
- **URL Parsing:** Extends [`net/url`](https://pkg.go.dev/net/url) parser to break down URLs into granular components.

## Installation

To install `hq-go-url`, run:

```bash
go get -v -u github.com/hueristiq/hq-go-url
```

Make sure your Go environment is set up properly (Go 1.x or later is recommended).

## Usage

### Extraction

The `extractor` package lets you scan text and pull out URLs using advanced regex patterns.

```go
package main

import (
    "fmt"
    "log"

    "github.com/hueristiq/hq-go-url/extractor"
)

func main() {
    e := extractor.New(extractor.WithScheme())

    regex := e.CompileRegex()

    text := "Check these out: ftp://ftp.example.com, https://secure.example.com, and mailto:someone@example.com."

    urls := regex.FindAllString(text, -1)

    fmt.Println("Extracted URLs:")

    for _, u := range urls {

        fmt.Println(u)
    }
}
```

You can customize how URLs are extracted by specifying URL schemes, hosts, or providing custom regular expression patterns.

- Extract URLs with Schemes Pattern:

	```go
	e := extractor.New(
		extractor.WithSchemePattern(`(?:https?|ftp)://`),
	)
	```

	This configuration will extract URLs with `http`, `https`, or `ftp` schemes.

- Extract URLs with Host Pattern:

	```go
	e := extractor.New(
		extractor.WithHostPattern(`(?:www\.)?example\.com`),
	)

	```

	This configuration will extract URLs that have hosts matching `www.example.com` or `example.com`.

### Parsing

The `parser` package extends Go's `net/url` package to include detailed domain breakdown.

```go
package main

import (
	"fmt"

	"github.com/hueristiq/hq-go-url/parser"
)

func main() {
	p := parser.New()

	parsed, err := p.Parse("https://subdomain.example.com:8080/path/file.txt")
	if err != nil {
		fmt.Println("Error parsing URL:", err)

		return
	}

    fmt.Printf("Scheme: %s\n", parsed.Scheme)
    fmt.Printf("Host: %s\n", parsed.Host)
    fmt.Printf("Hostname: %s\n", parsed.Hostname())
    fmt.Printf("Subdomain: %s\n", parsed.Domain.Subdomain)
    fmt.Printf("Second-Level Domain (SLD): %s\n", parsed.Domain.SecondLevelDomain)
    fmt.Printf("Top-Level Domain (TLD): %s\n", parsed.Domain.TopLevelDomain)
    fmt.Printf("Port: %s\n", parsed.Port())
    fmt.Printf("Path: %s\n", parsed.Path)
}
```

You can customize how URLs are parsed by specifying default scheme, or providing custom TLDs.

- Parse URLs with default scheme:

	```go
	p := parser.New(parser.WithDefaultScheme("https"))
	```

- Parse URLs with custom TLDs:

	```go
	p := parser.New(parser.WithTLDs("custom", "custom2"))
	```

## Contributing

Contributions are welcome and encouraged! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-url/pulls) or report [Issues](https://github.com/hueristiq/hq-go-url/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-url/blob/master/CONTRIBUTING.md).

A big thank you to all the [contributors](https://github.com/hueristiq/hq-go-url/graphs/contributors) for your ongoing support!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-url&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-url/blob/master/LICENSE).