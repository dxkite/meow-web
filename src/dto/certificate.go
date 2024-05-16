package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
)

// SSL证书
type Certificate struct {
	Id string `json:"id"`
	// 证书备注名
	Name string `json:"name"`
	// 证书支持的域名
	DNSNames []string `json:"dns_names"`
	// 证书开启时间
	NotBefore time.Time `json:"not_before"`
	// 证书有效期
	NotAfter time.Time `json:"not_after"`
	// 私钥
	Key string `json:"key,omitempty"`
	// 证书
	Certificate string    `json:"certificate,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewCertificate(item *entity.Certificate) *Certificate {
	obj := &Certificate{
		Id: identity.Format(constant.CertificatePrefix, item.Id),
	}
	obj.Name = item.Name
	obj.NotBefore = item.NotBefore
	obj.NotAfter = item.NotAfter
	obj.DNSNames = item.DNSNames
	obj.CreatedAt = item.CreatedAt
	obj.UpdatedAt = item.UpdatedAt
	return obj
}
