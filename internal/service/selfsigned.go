package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"os"
	"path"
	"time"
)

func LoadCertificate(base string) (*tls.Certificate, error) {
	crt, err := os.ReadFile(path.Join(base, "tls.crt"))
	if err != nil {
		return nil, err
	}
	key, err := os.ReadFile(path.Join(base, "tls.key"))
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func NewSelfSignedLocalhostCertificate() (*tls.Certificate, error) {
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		Issuer:                pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	return &tls.Certificate{
		Certificate: [][]byte{certificate},
		PrivateKey:  priv,
	}, nil
}
