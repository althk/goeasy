package grpcutils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type TLSConfig struct {
	CertFilePath     string
	KeyFilePath      string
	ClientCAFilePath string
	RootCAFilePath   string
	SkipTLS          bool
	NoClientCert     bool
}

func (c *TLSConfig) Creds() (credentials.TransportCredentials, error) {
	if c.SkipTLS {
		return insecure.NewCredentials(), nil
	}
	// init new tls config and load the cert
	cfg, err := c.newTLS()
	if err != nil {
		return nil, err
	}
	// if client ca is set, load it and enable client
	// verification
	if err = c.setClientCAs(cfg); err != nil {
		return nil, err
	}
	// if root ca is set, load it to enable server
	// verification
	if err = c.setRootCAs(cfg); err != nil {
		return nil, err
	}

	return credentials.NewTLS(cfg), nil
}

func (c *TLSConfig) newTLS() (*tls.Config, error) {
	cfg := &tls.Config{
		ClientAuth: tls.NoClientCert,
	}
	if c.NoClientCert {
		return cfg, nil
	}
	tlsKeyPair, err := tls.LoadX509KeyPair(c.CertFilePath, c.KeyFilePath)
	if err != nil {
		return nil, err
	}
	cfg.Certificates = []tls.Certificate{tlsKeyPair}
	return cfg, nil
}

func (c *TLSConfig) setRootCAs(cfg *tls.Config) error {
	if c.RootCAFilePath == "" {
		return nil
	}
	certPool, err := newCertPool(c.RootCAFilePath)
	if err != nil {
		return err
	}
	cfg.RootCAs = certPool

	return nil
}

func (c *TLSConfig) setClientCAs(cfg *tls.Config) error {
	if c.ClientCAFilePath == "" {
		return nil
	}
	certPool, err := newCertPool(c.ClientCAFilePath)
	if err != nil {
		return err
	}
	cfg.ClientCAs = certPool
	cfg.ClientAuth = tls.RequireAndVerifyClientCert

	return nil
}

// newCertPool creates a new CertPool and appends the cert
// at the given path to the pool.
func newCertPool(caPath string) (*x509.CertPool, error) {
	pemData, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemData) {
		return nil, fmt.Errorf("error adding CA cert to pool for %s", caPath)
	}
	return certPool, nil
}
