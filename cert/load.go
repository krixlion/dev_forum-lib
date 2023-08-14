package cert

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/krixlion/dev_forum-lib/fs"
)

// LoadX509KeyPair is just like tls.LoadX509KeyPair except is uses lib/fs
// instead of os package to interact with filesystem allowing to mock it out in tests.
func LoadX509KeyPair(certFilePath, keyFilePath string) (tls.Certificate, error) {
	keyPEMBlock, err := fs.ReadFile(keyFilePath)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEMBlock, err := fs.ReadFile(certFilePath)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certPEMBlock, keyPEMBlock)
}

// LoadCA loads CA certificate from given path and returns it or a non-nil error.
// The file can be mocked using lib/fs package.
func LoadCA(path string) (*x509.Certificate, error) {
	pemServerCA, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(pemServerCA)
}
