package tls

// func VerifyTLS() {
// 	// mTLS for one endpoint:
// 	// in your tls config don't demand a client cert and just verify it in the endpoint
// 	pr, ok := peer.FromContext(ctx)
// 	if !ok {
// 		return nil, errors.New("failed to get peer info")
// 	}

// 	tlsInfo, ok := pr.AuthInfo.(credentials.TLSInfo)
// 	if !ok {
// 		return nil, errors.New("unexpected peer transport credentials")
// 	}

// 	// Check the client certificate to make sure it's valid and issued by a trusted CA
// 	cert := tlsInfo.State.VerifiedChains[0][0]
// 	// verify the certificate is issued to the correct name
// 	if err := cert.VerifyHostname("first-service-name"); err != nil {
// 		return nil, err
// 	}
// }
