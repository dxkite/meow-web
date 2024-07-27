package monitor

type DynamicStat struct {
	Id             uint64  `gorm:"primarykey"`
	Time           uint64  `json:"time"`
	CpuPercent     float64 `json:"cpu_percent"`
	Load1          float64 `json:"load_1"`
	Load5          float64 `json:"load_5"`
	Load15         float64 `json:"load_15"`
	MemSwapUsed    uint64  `json:"mem_swap_used"`
	MemVirtualUsed uint64  `json:"mem_virtual_used"`
	NetRecv        uint64  `json:"net_recv"`
	NetSent        uint64  `json:"net_send"`
	DiskUsage      uint64  `json:"disk_usage"`
	DiskWrite      uint64  `json:"disk_write"`
	DiskRead       uint64  `json:"disk_read"`
}

func NewDynamicStat() *DynamicStat {
	entity := new(DynamicStat)
	return entity
}
