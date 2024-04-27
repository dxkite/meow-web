package entity

import (
	"crypto/tls"
	"crypto/x509"
	"time"
)

type Certificate struct {
	Base

	Name        string    `json:"name"`
	Domain      []string  `json:"domain" gorm:"serializer:json"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
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
	entity.StartTime = leaf.NotBefore
	entity.EndTime = leaf.NotAfter
	entity.Domain = leaf.DNSNames
	if len(entity.Domain) > 0 {
		entity.Name = entity.Domain[0]
	}
	return entity, nil
}
