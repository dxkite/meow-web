package entity

// 入口配置信息
type ServerName struct {
	Base
	Name          string `json:"string"`         // 主机名
	CertificateId uint64 `json:"certificate_id"` // 证书ID
}
