package entity

type Collection struct {
	Base

	// 树型节点部分
	ParentId uint64 `gorm:"index"`
	Index    string `gorm:"index"`
	Order    int
	Depth    int

	// 合辑名
	Name        string `gorm:"index"`
	Description string
}
