package cert

import (
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc/credentials"
)

func NewServerOptionalMTLSCreds(caCertPool *x509.CertPool, certs ...tls.Certificate) credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{
		Certificates: certs,
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    caCertPool,
	})
}

func NewClientMTLSCreds(caCertPool *x509.CertPool, certs ...tls.Certificate) credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{
		Certificates: certs,
		RootCAs:      caCertPool,
	})
}
