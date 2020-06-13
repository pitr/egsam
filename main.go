package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pitr/gig"
	"github.com/pitr/gig/middleware"
)

var port = "1965"

func main() {
	g := gig.New()

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "path=${uri} status=${status} duration=${latency} ${error}\n",
	}))
	g.Use(middleware.Recover())

	g.Static("", "static")

	v13 := g.Group("/v0.13.0")
	{
		v13.Handle("/1.1.write.timeout", func(c gig.Context) error {
			time.Sleep(10 * time.Minute)
			return c.Gemini(gig.StatusSuccess, "FAILED. You waited too long")
		})
		v13.Handle("/1.1.no.close", func(c gig.Context) error {
			c.Gemini(gig.StatusSuccess, "Nice, your client displayed partial content\n")
			time.Sleep(10 * time.Minute)
			_, err := c.Response().Write([]byte("FAILED"))
			return err
		})
		v13.Handle("/3.no.cr", func(c gig.Context) error {
			res := c.Response()
			res.Committed = true
			res.Status = gig.StatusSuccess
			res.Meta = "text/gemini"
			_, err := res.Write([]byte("20 text/gemini\nFail if you see this"))
			return err
		})
		v13.Handle("/3.1.bad.status", func(c gig.Context) error {
			res := c.Response()
			res.Committed = true
			res.Status = gig.StatusSuccess
			res.Meta = "text/gemini"
			_, err := res.Write([]byte("hi text/gemini\r\n"))
			return err
		})
		v13.Handle("/3.1.no.space", func(c gig.Context) error {
			res := c.Response()
			res.Committed = true
			res.Status = gig.StatusSuccess
			res.Meta = "text/gemini"
			_, err := res.Write([]byte("20text/gemini\r\n"))
			return err
		})
		v13.Handle("/3.1.long.meta", func(c gig.Context) error {
			return c.NoContent(gig.StatusSuccess, strings.Repeat(" ", 2000))
		})
		v13.Handle("/3.2.one.digit", func(c gig.Context) error {
			res := c.Response()
			res.Committed = true
			res.Status = gig.StatusSuccess
			res.Meta = "text/gemini"
			_, err := res.Write([]byte("2 text/gemini\r\nFAILED\n"))
			return err
		})
		v13.Handle("/3.2.three.digits", func(c gig.Context) error {
			res := c.Response()
			res.Committed = true
			res.Status = gig.StatusSuccess
			res.Meta = "text/gemini"
			_, err := res.Write([]byte("222 text/gemini\r\nFAILED\n"))
			return err
		})
		v13.Handle("/3.2.status.1", unknownStatus(19))
		v13.Handle("/3.2.status.2", unknownStatus(29))
		v13.Handle("/3.2.status.3", unknownStatus(39))
		v13.Handle("/3.2.status.4", unknownStatus(49))
		v13.Handle("/3.2.status.5", unknownStatus(59))
		v13.Handle("/3.2.status.6", unknownStatus(69))
		v13.Handle("/3.2.status.9", func(c gig.Context) error {
			res := c.Response()
			res.Committed = true
			res.Status = gig.StatusSuccess
			res.Meta = "text/gemini"
			_, err := res.Write([]byte("99 unknown status code\r\n"))
			return err
		})
		v13.Handle("/3.2.1.percent", func(c gig.Context) error {
			const expect = "99%20%2B%201%25"
			q := c.URL().RawQuery
			if q == "" {
				return c.NoContent(gig.StatusInput, "Please enter the following: 99 + 1%")
			}
			if q == expect {
				return c.Gemini(gig.StatusSuccess, "Passed\n=> 3.2.1.gmi Back")
			}
			return c.Gemini(gig.StatusSuccess, fmt.Sprintf("FAILED\nClient sent %v\nShould send %v\n=> 3.2.1.gmi Back", q, expect))
		})
		v13.Handle("/3.2.1.long", func(c gig.Context) error {
			q := c.URL().RawQuery
			if q == "" {
				return c.NoContent(gig.StatusInput, "Please enter the input as instructed on the previous page")
			}
			if q != strings.Repeat("x", len(q)) {
				return c.Gemini(gig.StatusSuccess, fmt.Sprintf("FAILED\nYour client sent magnled input\n```\n%v\n```\n=> 3.2.1.gmi Back", q))
			}
			return c.Gemini(gig.StatusSuccess, fmt.Sprintf("Your client sent %d bytes\n=> 3.2.1.gmi Back", len(q)))
		})
		v13.Handle("/3.2.2.text", func(c gig.Context) error {
			return c.Text(gig.StatusSuccess, "Pass")
		})
		v13.Handle("/3.2.2.html", func(c gig.Context) error {
			return c.Blob(gig.StatusSuccess, "text/html", []byte("<marquee>Pass</marquee>"))
		})
		v13.Handle("/3.2.2.jpg", func(c gig.Context) error {
			return c.File("static/pass.jpg")
		})
		v13.Handle("/3.2.2.jpg.bad", func(c gig.Context) error {
			return c.File("static/notimage.jpg")
		})
		v13.Handle("/3.2.3.redirect", func(c gig.Context) error {
			return c.NoContent(gig.StatusRedirectTemporary, "3.2.3.redirect.1")
		})
		v13.Handle("/3.2.3.redirect.1", func(c gig.Context) error {
			return c.NoContent(gig.StatusRedirectPermanent, "3.2.3.redirect.2")
		})
		v13.Handle("/3.2.3.redirect.2", func(c gig.Context) error {
			return c.Gemini(gig.StatusSuccess, "Pass\n=> 3.2.3.gmi Back")
		})
		v13.Handle("/3.2.4.fail", func(c gig.Context) error {
			return gig.NewErrorFrom(gig.ErrTemporaryFailure, "If you see this, client Passed")
		})
		v13.Handle("/3.2.5.fail", func(c gig.Context) error {
			return gig.NewErrorFrom(gig.ErrPermanentFailure, "If you see this, client Passed")
		})
		v13.Handle("/3.2.6.check", func(c gig.Context) error {
			return gig.NewErrorFrom(gig.ErrAuthorisedCertificateRequired, "If you see this, client Passed")
		})
		v13.Handle("/3.3.ascii", func(c gig.Context) error {
			return c.Blob(gig.StatusSuccess, "text/gemini; charset=us-ascii", []byte("Pass\n"))
		})
		v13.Handle("/3.3.utf8", func(c gig.Context) error {
			return c.Blob(gig.StatusSuccess, "text/gemini; charset=utf-8", []byte("𝒫𝒶𝓈𝓈"))
		})
		v13.Handle("/3.3.utf16", func(c gig.Context) error {
			return c.Blob(gig.StatusSuccess, "text/gemini; charset=utf-16", []byte("\x14\x6d\xc3\x15\x15\x24\xe2"))
		})
		v13.Handle("/3.3.utf8.bad", func(c gig.Context) error {
			return c.Blob(gig.StatusSuccess, "text/gemini; charset=utf-8", []byte("\x14\x6d\xc3\x15\x15\x24\xe2"))
		})
		v13.Handle("/3.4.text.unknown", func(c gig.Context) error {
			return c.Blob(gig.StatusSuccess, "text/garbage", []byte("Pass"))
		})
		v13.Handle("/4.3.transient", func(c gig.Context) error {
			if c.Certificate() == nil {
				return gig.ErrTransientCertificateRequested
			}
			return c.Gemini(gig.StatusSuccessEndOfClientCertificateSession, `## Almost there

Now clicking below should ONLY offer to create a new certificate. Please create it.

=> 4.3.transient.2 Continue`)
		})
		v13.Handle("/4.3.transient.2", func(c gig.Context) error {
			if c.Certificate() == nil {
				return gig.NewErrorFrom(gig.ErrTransientCertificateRequested, "Please create a new transient certificate")
			}
			return c.Gemini(gig.StatusSuccessEndOfClientCertificateSession, `Pass if client created a new certificate.
Fail if client offered to use another transient certificate, or did not present user with choice.

=> 4.3.gmi Back`)
		})
	}

	panic(g.Run(":"+port, "egsam.crt", "egsam.key"))
}

func unknownStatus(i int) gig.HandlerFunc {
	return func(c gig.Context) error {
		if c.URL().RawQuery != "" {
			return c.Gemini(gig.StatusSuccess, "Pass")
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
