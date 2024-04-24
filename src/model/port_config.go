package model

import "gorm.io/gorm"

// 入口配置信息
type PortConfig struct {
	gorm.Model
	// 监听域名端口
	// hostname:433
	Listen []string `gorm:"serializer:json" json:"listen"`
	// 校验类型
	Authorization string `json:"authorization"`
}
