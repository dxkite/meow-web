package dto

import (
	"fmt"
	"sort"
	"strconv"

	"dxkite.cn/meownest/src/entity"
)

type DynamicStatCollection struct {
	Time           []uint64  `json:"time"`
	CpuPercent     []float64 `json:"cpu_percent"`
	Load1          []float64 `json:"load_1"`
	Load5          []float64 `json:"load_5"`
	Load15         []float64 `json:"load_15"`
	MemSwapUsed    []float64 `json:"mem_swap_used"`
	MemVirtualUsed []float64 `json:"mem_virtual_used"`
	NetRecv        []uint64  `json:"net_recv"`
	NetRecvSpeed   []float64 `json:"net_recv_speed"`
	NetSent        []uint64  `json:"net_send"`
	NetSentSpeed   []float64 `json:"net_send_speed"`
	DiskUsage      []float64 `json:"disk_usage"`
	DiskWrite      []uint64  `json:"disk_write"`
	DiskWriteSpeed []float64 `json:"disk_write_speed"`
	DiskRead       []uint64  `json:"disk_read"`
	DiskReadSpeed  []float64 `json:"disk_read_speed"`
}

func NewDynamicStatCollection(entities []*entity.DynamicStat, msTotal, mvTotal, dTotal uint64) *DynamicStatCollection {
	coll := &DynamicStatCollection{}
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Time < entities[j].Time
	})

	coll.NetRecvSpeed = append(coll.NetRecvSpeed, 0)
	coll.NetSentSpeed = append(coll.NetSentSpeed, 0)
	coll.DiskWriteSpeed = append(coll.DiskWriteSpeed, 0)
	coll.DiskReadSpeed = append(coll.DiskReadSpeed, 0)

	for i, v := range entities {
		coll.Time = append(coll.Time, v.Time)
		coll.CpuPercent = append(coll.CpuPercent, formatFloat64(v.CpuPercent))
		coll.Load1 = append(coll.Load1, formatFloat64(v.Load1))
		coll.Load5 = append(coll.Load5, formatFloat64(v.Load5))
		coll.Load15 = append(coll.Load15, formatFloat64(v.Load15))
		coll.MemSwapUsed = append(coll.MemSwapUsed, formatFloat64(float64(v.MemSwapUsed)/float64(msTotal)*100))
		coll.MemVirtualUsed = append(coll.MemVirtualUsed, formatFloat64(float64(v.MemVirtualUsed)/float64(mvTotal)*100))
		coll.NetRecv = append(coll.NetRecv, v.NetRecv)
		coll.NetSent = append(coll.NetSent, v.NetSent)
		coll.DiskUsage = append(coll.DiskUsage, formatFloat64(float64(v.DiskUsage)/float64(dTotal)*100))
		coll.DiskWrite = append(coll.DiskWrite, v.DiskWrite)
		coll.DiskRead = append(coll.DiskRead, v.DiskRead)

		if i > 0 {
			prev := entities[i-1]
			gap := v.Time - prev.Time

			coll.NetRecvSpeed = append(coll.NetRecvSpeed, formatFloat64(float64(v.NetRecv-prev.NetRecv)/float64(gap)))
			coll.NetSentSpeed = append(coll.NetSentSpeed, formatFloat64(float64(v.NetSent-prev.NetSent)/float64(gap)))
			coll.DiskWriteSpeed = append(coll.DiskWriteSpeed, formatFloat64(float64(v.DiskWrite-prev.DiskWrite)/float64(gap)))
			coll.DiskReadSpeed = append(coll.DiskReadSpeed, formatFloat64(float64(v.DiskRead-prev.DiskRead)/float64(gap)))
		}
	}
	return coll
}

func formatFloat64(v float64) float64 {
	vv, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", v), 64)
	return vv
}
