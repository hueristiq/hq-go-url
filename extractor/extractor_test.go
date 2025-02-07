package extractor_test

import (
	"testing"

	"go.source.hueristiq.com/url/extractor"
)

func TestNewExtractor(t *testing.T) {
	t.Parallel()

	e := extractor.New()

	if e == nil {
		t.Error("New() = nil; want non-nil")
	}
}

func TestCompileRegex(t *testing.T) {
	t.Parallel()

	e := extractor.New()

	regex := e.CompileRegex()

	if regex == nil {
		t.Errorf("CompileRegex() = nil; want non-nil")
	}
}

func TestExtractorWithScheme(t *testing.T) {
	t.Parallel()

	e := extractor.New(
		extractor.WithScheme(),
	)

	regex := e.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			http://localhost
			http://localhost:8000/home

			https://www.example.com
			http://www.example.com:8080/page
			https://www.example.com/search?q=openai
			https://www.example.com/resources/docs/guide
			https://www.example.com/about#team
			https://user:password@example.com
			https://www.example.com:8080/search?q=openai#results

			https://www.例子.公司.cn
			http://中国.中国/中国
			http://中国.中国/foo中国

			http://192.168.1.1/setup
			http://192.168.1.1:8080/setup
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080

			https://example.com/products?id=50&category=books
			https://example.com/search?name=John%20Doe&age=30

			http://foo.com/🐼
			https://shmibbles.me/tmp/自殺でも？.png
			https://subdomain.example.co.uk

			foo://example.com
			file:///C:/path/to/file.txt
			mailto:user@example.com
			ftp://user:password@ftp.example.com:21
			data:text/plain;base64,aGVsbG8=
			javascript:alert('Hello World');
			sms:123
			xmpp:foo@bar
			bitcoin:Addr23?amount=1&message=foo
			cid:foo-32x32.v2_fe0f1423.png
			mid:960830.1639@XIson.com
			postgres://user:pass@host.com:5432/path?k=v#f
			zoommtg://zoom.us/join?confno=1234&pwd=xxx

			https://www.example.com/test space
			https://www.example.com, it's cool!

			/path/to/resource
			www.example.com
			www.example.com/resources
			www.example.com/search?q=openai
			192.168.1.1
			192.168.1.1:8080
			::1
			::ffff:0:0
			64:ff9b::
			64:ff9b:1::
			100::
			2001::
			2001:1::1
			2001:1::2
			2001:2::
			2001:3::
			2001:4:112::
			2001:10::
			2001:20::
			2002::
			2620:4f:8000::
			fc00::
			fe80::
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"http://localhost",
				"http://localhost:8000/home",
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://www.例子.公司.cn",
				"http://中国.中国/中国",
				"http://中国.中国/foo中国",
				// "http://उदाहरण.परीकषा",
				"http://192.168.1.1/setup",
				"http://192.168.1.1:8080/setup",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				// "http://✪foo.bar/pa✪th©more",
				"http://foo.com/🐼",
				"https://shmibbles.me/tmp/自殺でも？.png",
				"https://subdomain.example.co.uk",
				"foo://example.com",
				"file:///C:/path/to/file.txt",
				"mailto:user@example.com",
				"ftp://user:password@ftp.example.com:21",
				"sms:123",
				"xmpp:foo@bar",
				"bitcoin:Addr23?amount=1&message=foo",
				"cid:foo-32x32.v2_fe0f1423.png",
				"mid:960830.1639@XIson.com",
				"postgres://user:pass@host.com:5432/path?k=v#f",
				"zoommtg://zoom.us/join?confno=1234&pwd=xxx",
				"https://www.example.com/test",
				"https://www.example.com",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

func TestExtractorWithSchemePattern(t *testing.T) {
	t.Parallel()

	e := extractor.New(
		extractor.WithSchemePattern(`(?:https?://)`),
	)

	regex := e.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			http://localhost
			http://localhost:8000/home

			https://www.example.com
			http://www.example.com:8080/page
			https://www.example.com/search?q=openai
			https://www.example.com/resources/docs/guide
			https://www.example.com/about#team
			https://user:password@example.com
			https://www.example.com:8080/search?q=openai#results

			https://www.例子.公司.cn
			http://中国.中国/中国
			http://中国.中国/foo中国

			http://192.168.1.1/setup
			http://192.168.1.1:8080/setup
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080

			https://example.com/products?id=50&category=books
			https://example.com/search?name=John%20Doe&age=30
			http://foo.com/🐼
			https://shmibbles.me/tmp/自殺でも？.png
			https://subdomain.example.co.uk

			foo://example.com
			file:///C:/path/to/file.txt
			mailto:user@example.com
			ftp://user:password@ftp.example.com:21
			data:text/plain;base64,aGVsbG8=
			javascript:alert('Hello World');
			sms:123
			xmpp:foo@bar
			bitcoin:Addr23?amount=1&message=foo
			cid:foo-32x32.v2_fe0f1423.png
			mid:960830.1639@XIson.com
			postgres://user:pass@host.com:5432/path?k=v#f
			zoommtg://zoom.us/join?confno=1234&pwd=xxx

			https://www.example.com/test space
			https://www.example.com, it's cool!

			/path/to/resource
			www.example.com
			www.example.com/resources
			www.example.com/search?q=openai
			192.168.1.1
			192.168.1.1:8080
			::1
			::ffff:0:0
			64:ff9b::
			64:ff9b:1::
			100::
			2001::
			2001:1::1
			2001:1::2
			2001:2::
			2001:3::
			2001:4:112::
			2001:10::
			2001:20::
			2002::
			2620:4f:8000::
			fc00::
			fe80::
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"http://localhost",
				"http://localhost:8000/home",
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://www.例子.公司.cn",
				"http://中国.中国/中国",
				"http://中国.中国/foo中国",
				// "http://उदाहरण.परीकषा",
				"http://192.168.1.1/setup",
				"http://192.168.1.1:8080/setup",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				// "http://✪foo.bar/pa✪th©more",
				"http://foo.com/🐼",
				"https://shmibbles.me/tmp/自殺でも？.png",
				"https://subdomain.example.co.uk",
				"https://www.example.com/test",
				"https://www.example.com",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

func TestExtractorWithHost(t *testing.T) {
	t.Parallel()

	e := extractor.New(
		extractor.WithHost(),
	)

	regex := e.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			http://localhost
			http://localhost:8000/home

			https://www.example.com
			http://www.example.com:8080/page
			https://www.example.com/search?q=openai
			https://www.example.com/resources/docs/guide
			https://www.example.com/about#team
			https://user:password@example.com
			https://www.example.com:8080/search?q=openai#results

			https://www.例子.公司.cn
			http://中国.中国/中国
			http://中国.中国/foo中国
			http://उदाहरण.परीकषा

			http://192.168.1.1/setup
			http://192.168.1.1:8080/setup
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080

			https://example.com/products?id=50&category=books
			https://example.com/search?name=John%20Doe&age=30
			http://✪foo.bar/pa✪th©more
			http://foo.com/🐼
			https://shmibbles.me/tmp/自殺でも？.png
			https://subdomain.example.co.uk

			foo://example.com
			file:///C:/path/to/file.txt
			mailto:user@example.com
			ftp://user:password@ftp.example.com:21
			data:text/plain;base64,aGVsbG8=
			javascript:alert('Hello World');
			sms:123
			xmpp:foo@bar
			bitcoin:Addr23?amount=1&message=foo
			cid:foo-32x32.v2_fe0f1423.png
			mid:960830.1639@XIson.com
			postgres://user:pass@host.com:5432/path?k=v#f
			zoommtg://zoom.us/join?confno=1234&pwd=xxx

			https://www.example.com/test space
			https://www.example.com, it's cool!

			/path/to/resource
			www.example.com
			www.example.com/resources
			www.example.com/search?q=openai
			192.168.1.1
			192.168.1.1:8080
			::1
			::ffff:0:0
			64:ff9b::
			64:ff9b:1::
			100::
			2001::
			2001:1::1
			2001:1::2
			2001:2::
			2001:3::
			2001:4:112::
			2001:10::
			2001:20::
			2002::
			2620:4f:8000::
			fc00::
			fe80::
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"http://localhost",
				"http://localhost:8000/home",
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://www.例子.公司.cn",
				"http://中国.中国/中国",
				"http://中国.中国/foo中国",
				"http://उदाहरण.परीकषा",
				"http://192.168.1.1/setup",
				"http://192.168.1.1:8080/setup",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				"http://✪foo.bar/pa✪th©more",
				"http://foo.com/🐼",
				"https://shmibbles.me/tmp/自殺でも？.png",
				"https://subdomain.example.co.uk",
				"foo://example.com",
				"file:///C:/path/to/file.txt",
				"mailto:user@example.com",
				"ftp://user:password@ftp.example.com:21",
				"sms:123",
				"xmpp:foo@bar",
				"bitcoin:Addr23?amount=1&message=foo",
				"cid:foo-32x32.v2_fe0f1423.png",
				"mid:960830.1639@XIson.com",
				"postgres://user:pass@host.com:5432/path?k=v#f",
				"zoommtg://zoom.us/join?confno=1234&pwd=xxx",
				"https://www.example.com/test",
				"https://www.example.com",
				"www.example.com",
				"www.example.com/resources",
				"www.example.com/search?q=openai",
				"192.168.1.1",
				"192.168.1.1:8080",
				"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

func TestExtractorWithHostPattern(t *testing.T) {
	t.Parallel()

	e := extractor.New(
		extractor.WithHostPattern(`(?:(?:\w+[.])*example\.com` + extractor.ExtractorPortOptionalPattern + `)`),
	)

	regex := e.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			http://localhost
			http://localhost:8000/home

			https://www.example.com
			http://www.example.com:8080/page
			https://www.example.com/search?q=openai
			https://www.example.com/resources/docs/guide
			https://www.example.com/about#team
			https://user:password@example.com
			https://www.example.com:8080/search?q=openai#results

			https://www.例子.公司.cn
			http://中国.中国/中国
			http://中国.中国/foo中国
			http://उदाहरण.परीकषा

			http://192.168.1.1/setup
			http://192.168.1.1:8080/setup
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080

			https://example.com/products?id=50&category=books
			https://example.com/search?name=John%20Doe&age=30
			http://✪foo.bar/pa✪th©more
			http://foo.com/🐼
			https://shmibbles.me/tmp/自殺でも？.png
			https://subdomain.example.co.uk

			foo://example.com
			file:///C:/path/to/file.txt
			mailto:user@example.com
			ftp://user:password@ftp.example.com:21
			data:text/plain;base64,aGVsbG8=
			javascript:alert('Hello World');
			sms:123
			xmpp:foo@bar
			bitcoin:Addr23?amount=1&message=foo
			cid:foo-32x32.v2_fe0f1423.png
			mid:960830.1639@XIson.com
			postgres://user:pass@host.com:5432/path?k=v#f
			zoommtg://zoom.us/join?confno=1234&pwd=xxx

			https://www.example.com/test space
			https://www.example.com, it's cool!

			/path/to/resource
			www.example.com
			www.example.com/resources
			www.example.com/search?q=openai
			192.168.1.1
			192.168.1.1:8080
			::1
			::ffff:0:0
			64:ff9b::
			64:ff9b:1::
			100::
			2001::
			2001:1::1
			2001:1::2
			2001:2::
			2001:3::
			2001:4:112::
			2001:10::
			2001:20::
			2002::
			2620:4f:8000::
			fc00::
			fe80::
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				"foo://example.com",
				"mailto:user@example.com",
				"ftp://user:password@ftp.example.com:21",
				"https://www.example.com/test",
				"https://www.example.com",
				"www.example.com",
				"www.example.com/resources",
				"www.example.com/search?q=openai",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

// equalSlices checks if two slices of strings are equal.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
