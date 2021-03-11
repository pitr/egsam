package main

import (
	"crypto/tls"
	"io/ioutil"
	"strings"

	"github.com/pitr/gig"
)

func main() {
	g := gig.Default()

	g.TLSConfig.MinVersion = tls.VersionTLS13

	g.TLSConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if !strings.Contains(hello.ServerName, ".glv.one") {
			return nil, nil
		}

		c, err := ioutil.ReadFile("/meta/credentials/letsencrypt/current/fullchain.pem")
		if err != nil {
			return nil, err
		}
		k, err := ioutil.ReadFile("/meta/credentials/letsencrypt/current/privkey.pem")
		if err != nil {
			return nil, err
		}
		cert, err := tls.X509KeyPair(c, k)
		if err != nil {
			return nil, err
		}
		return &cert, nil
	}

	g.Handle("/", func(c gig.Context) error {
		return c.Gemini("Passed, your client supports servers requiring TLS v1.3\n=> gemini://egsam.glv.one/4.1.gmi Back")
	})

	panic(g.Run("egsam.crt", "egsam.key"))
}
