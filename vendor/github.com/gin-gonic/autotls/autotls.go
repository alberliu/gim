package autotls

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

// Run support 1-line LetsEncrypt HTTPS servers
func Run(r http.Handler, domain ...string) error {
	return http.Serve(autocert.NewListener(domain...), r)
}

// RunWithManager support custom autocert manager
func RunWithManager(r http.Handler, m *autocert.Manager) error {
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   r,
	}

	go http.ListenAndServe(":http", m.HTTPHandler(nil))

	return s.ListenAndServeTLS("", "")
}
