package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func ParsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no rsa key")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return pri, nil
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}

func ParsePublicKeyFromCertificate(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no rsa key")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	k := cert.PublicKey
	switch t := k.(type) {
	case *rsa.PublicKey:
		return t, nil
	}
	return nil, fmt.Errorf("unsupported key type %T", k)
}
