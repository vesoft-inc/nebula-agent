package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func LoadTLSConfig(caCert, cert, key []byte) (*tls.Config, error) {
	c, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	rootCAPool := x509.NewCertPool()
	ok := rootCAPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("failed to load cert files")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{c},
		RootCAs:      rootCAPool,
	}, nil
}

func GetCerts(caPath, certPath, keyPath string) ([]byte, []byte, []byte, error) {
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, nil, nil, err
	}
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, nil, err
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, nil, err
	}
	return caCert, cert, key, nil
}
