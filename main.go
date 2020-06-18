package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pitr/gig"
)

func main() {
	g := gig.Default()

	g.Static("", "static")

	g.Handle("/1.1.write.timeout", func(c gig.Context) error {
		time.Sleep(10 * time.Minute)
		return c.Gemini("FAILED. You waited too long")
	})
	g.Handle("/1.1.no.close", func(c gig.Context) error {
		c.Gemini("Nice, your client displayed partial content\n")
		time.Sleep(10 * time.Minute)
		_, err := c.Response().Write([]byte("FAILED"))
		return err
	})
	g.Handle("/3.no.cr", func(c gig.Context) error {
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		_, err := res.Write([]byte("20 text/gemini\nFail if you see this"))
		return err
	})
	g.Handle("/3.1.bad.status", func(c gig.Context) error {
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		_, err := res.Write([]byte("hi text/gemini\r\n"))
		return err
	})
	g.Handle("/3.1.no.space", func(c gig.Context) error {
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		_, err := res.Write([]byte("20text/gemini\r\n"))
		return err
	})
	g.Handle("/3.1.long.meta", func(c gig.Context) error {
		return c.NoContent(gig.StatusSuccess, strings.Repeat(" ", 2000))
	})
	g.Handle("/3.2.one.digit", func(c gig.Context) error {
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		_, err := res.Write([]byte("2 text/gemini\r\nFAILED\n"))
		return err
	})
	g.Handle("/3.2.three.digits", func(c gig.Context) error {
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		_, err := res.Write([]byte("222 text/gemini\r\nFAILED\n"))
		return err
	})
	g.Handle("/3.2.status.1", unknownStatus(19))
	g.Handle("/3.2.status.2", unknownStatus(29))
	g.Handle("/3.2.status.3", unknownStatus(39))
	g.Handle("/3.2.status.4", unknownStatus(49))
	g.Handle("/3.2.status.5", unknownStatus(58))
	g.Handle("/3.2.status.6", unknownStatus(69))
	g.Handle("/3.2.status.9", func(c gig.Context) error {
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		_, err := res.Write([]byte("99 unknown status code\r\n"))
		return err
	})
	g.Handle("/3.2.1.percent", func(c gig.Context) error {
		const expect = "99%20%2B%201%25"
		q := c.URL().RawQuery
		if q == "" {
			return c.NoContent(gig.StatusInput, "Please enter the following: 99 + 1%")
		}
		if q == expect {
			return c.Gemini("Passed\n=> 3.2.1.gmi Back")
		}
		return c.Gemini(`# FAILED
Client sent %v
Should send %v
=> 3.2.1.gmi Back`, q, expect)
	})
	g.Handle("/3.2.1.long", func(c gig.Context) error {
		q := c.URL().RawQuery
		if q == "" {
			return c.NoContent(gig.StatusInput, "Please enter the input as instructed on the previous page")
		}
		if q != strings.Repeat("x", len(q)) {
			return c.Gemini("FAILED\nYour client sent magnled input\n```\n%v\n```\n=> 3.2.1.gmi Back", q)
		}
		return c.Gemini("Your client sent %d bytes\n=> 3.2.1.gmi Back", len(q))
	})
	g.Handle("/3.2.2.text", func(c gig.Context) error {
		return c.Text("Pass")
	})
	g.Handle("/3.2.2.html", func(c gig.Context) error {
		return c.Blob("text/html", []byte("<marquee>Pass</marquee>"))
	})
	g.Handle("/3.2.2.jpg", func(c gig.Context) error {
		return c.File("static/pass.jpg")
	})
	g.Handle("/3.2.2.jpg.bad", func(c gig.Context) error {
		return c.File("static/notimage.jpg")
	})
	g.Handle("/3.2.3.redirect", func(c gig.Context) error {
		return c.NoContent(gig.StatusRedirectTemporary, "3.2.3.redirect.1")
	})
	g.Handle("/3.2.3.redirect.1", func(c gig.Context) error {
		return c.NoContent(gig.StatusRedirectPermanent, "3.2.3.redirect.2")
	})
	g.Handle("/3.2.3.redirect.2", func(c gig.Context) error {
		return c.Gemini("Pass\n=> 3.2.3.gmi Back")
	})
	g.Handle("/3.2.4.fail", func(c gig.Context) error {
		return gig.NewErrorFrom(gig.ErrTemporaryFailure, "If you see this, client Passed")
	})
	g.Handle("/3.2.5.fail", func(c gig.Context) error {
		return gig.NewErrorFrom(gig.ErrPermanentFailure, "If you see this, client Passed")
	})
	g.Handle("/3.2.6.check", func(c gig.Context) error {
		return gig.NewErrorFrom(gig.ErrCertificateNotAuthorised, "If you see this, client Passed")
	})
	g.Handle("/3.3.ascii", func(c gig.Context) error {
		return c.Blob("text/gemini; charset=us-ascii", []byte("Pass\n"))
	})
	g.Handle("/3.3.utf8", func(c gig.Context) error {
		return c.Blob("text/gemini; charset=utf-8", []byte("𝒫𝒶𝓈𝓈"))
	})
	g.Handle("/3.3.utf16", func(c gig.Context) error {
		return c.Blob("text/gemini; charset=utf-16", []byte("\x14\x6d\xc3\x15\x15\x24\xe2"))
	})
	g.Handle("/3.3.utf8.bad", func(c gig.Context) error {
		return c.Blob("text/gemini; charset=utf-8", []byte("\x14\x6d\xc3\x15\x15\x24\xe2"))
	})
	g.Handle("/3.4.text.unknown", func(c gig.Context) error {
		return c.Blob("text/garbage", []byte("Pass"))
	})
	g.Handle("/4.3.cert", func(c gig.Context) error {
		if c.Certificate() == nil {
			return gig.ErrClientCertificateRequired
		}
		name := c.Certificate().Subject.CommonName
		if name != "tester" {
			return c.Gemini("# FAILED\nSubject name '%s' does not match 'tester'", name)
		}
		return c.NoContent(gig.StatusRedirectTemporary, "4.3.cert.2")
	})
	g.Handle("/4.3.cert.2", func(c gig.Context) error {
		if c.Certificate() != nil {
			return c.Gemini(`## Almost there
Now deactivate current certificate (and delete it). Then refresh this page.`)
		}
		return c.Gemini(`## Pass
=> 4.3.gmi Back`)
	})

	panic(g.Run("egsam.crt", "egsam.key"))
}

func unknownStatus(i int) gig.HandlerFunc {
	return func(c gig.Context) error {
		if c.URL().RawQuery != "" {
			return c.Gemini("Pass")
		}
		res := c.Response()
		res.Committed = true
		res.Status = gig.StatusSuccess
		res.Meta = "text/gemini"
		switch i {
		case 19:
			_, err := res.Write([]byte(fmt.Sprintf("%d just enter anything\r\n", i)))
			return err
		case 29:
			_, err := res.Write([]byte("29 text/gemini\r\nGreat, your client handles future status codes\n"))
			return err
		case 39:
			_, err := res.Write([]byte("39 3.2.gmi\r\n"))
			return err
		default:
			_, err := res.Write([]byte(fmt.Sprintf("%d unknown status code\r\n", i)))
			return err
		}
	}
}
