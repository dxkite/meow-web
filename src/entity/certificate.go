package entity

import (
	"crypto/tls"
	"crypto/x509"
	"time"
)

type Certificate struct {
	Base

	Name        string    `json:"name"`
	DNSNames    []string  `json:"dns_names" gorm:"serializer:json"`
	NotBefore   time.Time `json:"not_before"`
	NotAfter    time.Time `json:"not_after"`
	Key         string    `json:"key"`
	Certificate string    `json:"certificate"`
}

func NewCertificateWithCertificateKey(certStr, keyStr string) (*Certificate, error) {
	cert, err := tls.X509KeyPair([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, err
	}

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	entity := &Certificate{}
	entity.Key = keyStr
	entity.Certificate = certStr
	entity.NotBefore = leaf.NotBefore
	entity.NotAfter = leaf.NotAfter
	entity.DNSNames = leaf.DNSNames
	return entity, nil
}
