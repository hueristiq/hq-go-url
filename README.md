# hq-go-url

[![go report card](https://goreportcard.com/badge/github.com/hueristiq/hq-go-url)](https://goreportcard.com/report/github.com/hueristiq/hq-go-url) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-url.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-url/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-url.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-url/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-url/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-url/blob/master/CONTRIBUTING.md)

`hq-go-url` is a [Go (Golang)](http://golang.org/) package designed for extracting, parsing, and manipulating URLs with ease. This library is useful for developers who need to work with URLs in a structured way.

## Resources

* [Features](#features)
* [Usage](#usage)
    * [URL Extraction](#url-extraction)
        * [Customizing URL Extraction](#customizing-url-extraction)
    * [Domain Parsing](#domain-parsingn)
    * [URL Parsing](#url-parsing)
* [Contributing](#contributing)
* [Licensing](#licensing)
* [Credits](#credits)
    * [Contributors](#contributors)
    * [Similar Projects](#similar-projects)

## Features

* **Flexible URL Extraction:** Extract URLs from text using regular expressions.
* **Domain Parsing:** Parse domains into subdomains, root domains, and top-level domains (TLDs).
* **Extended URL Parsing:** Extend the standard `net/url` package in Go with additional fields and capabilities.

## Installation

To install the package, run the following command in your terminal:

```bash
go get -v -u github.com/hueristiq/hq-go-url
```

This command will download and install the `hq-go-url` package into your Go workspace, making it available for use in your projects.

## Usage

Below are examples demonstrating how to use the different features of the `hq-go-url` package.

### URL Extraction

You can extract URLs from a given text string using the Extractor. Here's a simple example:

```go
package main

import (
    "fmt"
    hqgourl "github.com/hueristiq/hq-go-url"
    "regexp"
)

func main() {
    extr := hqgourl.NewExtractor()
    text := "Check out this website: https://example.com and send an email to info@example.com."
    
    regex := extr.CompileRegex()
    matches := regex.FindAllString(text, -1)
    
    fmt.Println("Found URLs:", matches)
}
```

#### Customizing URL Extraction

You can customize how URLs are extracted by specifying URL schemes, hosts, or providing custom regular expression patterns.

* Extract URLs with Specific Schemes (e.g., HTTP, HTTPS, FTP):

    ```go
    extr := hqgourl.NewExtractor(
        hqgourl.ExtractorWithSchemePattern(`(?:https?|ftp)://`),
    )
    ```

    This configuration will extract only URLs starting with http, https, or ftp schemes.

* Extract URLs with Custom Host Patterns (e.g., example.com):

    ```go
    extr := hqgourl.NewExtractor(
        hqgourl.ExtractorWithHostPattern(`(?:www\.)?example\.com`),
    )

    ```

    This setup will extract URLs that have hosts matching www.example.com or example.com.

> [!NOTE]
> Since API is centered around [regexp.Regexp](https://golang.org/pkg/regexp/#Regexp), many other methods are available

### Domain Parsing

The `DomainParser` can parse domains into their components, such as subdomains, root domains, and TLDs:

```go
package main

import (
    "fmt"
    hqgourl "github.com/hueristiq/hq-go-url"
)

func main() {
    dp := hqgourl.NewDomainParser()

    parsedDomain := dp.Parse("subdomain.example.com")

    fmt.Printf("Subdomain: %s, Root Domain: %s, TLD: %s\n", parsedDomain.Sub, parsedDomain.Root, parsedDomain.TopLevel)
}
```

### URL Parsing

The `Parser` provides an extended way to parse URLs, including additional fields like port and file extension:

```go
package main

import (
    "fmt"
    hqgourl "github.com/hueristiq/hq-go-url"
)

func main() {
    up := hqgourl.NewParser()

    parsedURL, err := up.Parse("https://subdomain.example.com:8080/path/file.txt")
    if err != nil {
        fmt.Println("Error parsing URL:", err)

        return
    }

    fmt.Printf("Subdomain: %s\n", parsedURL.Domain.Sub)
    fmt.Printf("Root Domain: %s\n", parsedURL.Domain.Root)
    fmt.Printf("TLD: %s\n", parsedURL.Domain.TopLevel)
    fmt.Printf("Port: %d\n", parsedURL.Port)
    fmt.Printf("File Extension: %s\n", parsedURL.Extension)
}
```

Set a default scheme:

```go
up := hqgourl.NewParser(hqgourl.ParserWithDefaultScheme("https"))
```

## Contributing

We welcome contributions! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-url/pulls) or report [Issues](https://github.com/hueristiq/hq-go-url/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-url/blob/master/CONTRIBUTING.md).


## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-url/blob/master/LICENSE).

## Credits

### Contributors

A huge thanks to all the contributors who have helped make `hq-go-url` what it is today!

[![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-url&max=500)](https://github.com/hueristiq/hq-go-url/graphs/contributors)

### Similar Projects

If you're interested in more packages like this, check out:

[DomainParser](https://github.com/Cgboal/DomainParser) ◇ [urlx](https://github.com/goware/urlx) ◇ [xurls](https://github.com/mvdan/xurls) ◇ [goware's tldomains](https://github.com/goware/tldomains) ◇ [jakewarren's tldomains](https://github.com/jakewarren/tldomains)