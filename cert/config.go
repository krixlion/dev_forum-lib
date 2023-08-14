package cert

import (
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc/credentials"
)

func ServerOptionalMTLSCreds(caCertPool *x509.CertPool, certs ...tls.Certificate) credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{
		Certificates: certs,
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    caCertPool,
	})
}

func ClientMTLSCreds(caCertPool *x509.CertPool, certs ...tls.Certificate) credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{
		Certificates: certs,
		RootCAs:      caCertPool,
	})
}
