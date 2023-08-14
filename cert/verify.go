package cert

import (
	"context"
	"errors"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// VerifyClientTLS verifies client's certificate and returns a non-nil error
// describing the issue or nil if the cert is valid
//
// This function checks only if client's cert was provided,
// and whether DNSNames contain given hostname.
//
// Certificate is only propagated through context if grpc.Creds() is used.
func VerifyClientTLS(ctx context.Context, hostname string) error {
	tlsInfo, err := infoFromContext(ctx)
	if err != nil {
		return err
	}

	if len(tlsInfo.State.VerifiedChains) < 1 {
		return errors.New("client certificate not found")
	}

	cert := tlsInfo.State.VerifiedChains[0][0]

	return cert.VerifyHostname(hostname)
}

// infoFromContext returns TLS credentials extracted from handler's context or a non-nil err.
// Credentials are only propagated through context if grpc.Creds() is used.
func infoFromContext(ctx context.Context) (credentials.TLSInfo, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return credentials.TLSInfo{}, errors.New("failed to get peer info")
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return credentials.TLSInfo{}, errors.New("unexpected peer transport credentials")
	}

	return tlsInfo, nil
}
