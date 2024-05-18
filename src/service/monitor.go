package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"dxkite.cn/meownest/pkg/stat"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/value"
)

type Monitor interface {
	Collection(ctx context.Context) error
	GetRealtimeStat(ctx context.Context) (*GetRealtimeLoadStatResult, error)
}

type monitor struct {
	interval        int
	maxInterval     int
	memSwapTotal    uint64
	memVirtualTotal uint64
	diskTotal       uint64
	status          []*value.DynamicStat
}

func NewMonitor(interval, maxInterval int) Monitor {
	return &monitor{interval: interval, maxInterval: maxInterval}
}

type GetRealtimeLoadStatResult struct {
	Collection      *dto.LoadStatCollection `json:"collection"`
	MemSwapTotal    uint64                  `json:"mem_swap_total"`
	MemVirtualTotal uint64                  `json:"mem_virtual_total"`
	DiskTotal       uint64                  `json:"disk_total"`
}

func (m monitor) GetRealtimeStat(ctx context.Context) (*GetRealtimeLoadStatResult, error) {
	coll := &dto.LoadStatCollection{}

	for _, v := range m.status {
		coll.Time = append(coll.Time, v.Time)
		coll.CpuPercent = append(coll.CpuPercent, v.CpuPercent)
		coll.Load1 = append(coll.Load1, v.Load1)
		coll.Load5 = append(coll.Load5, v.Load5)
		coll.Load15 = append(coll.Load15, v.Load15)
		coll.MemSwapUsed = append(coll.MemSwapUsed, v.MemSwapUsed)
		coll.MemVirtualUsed = append(coll.MemVirtualUsed, v.MemVirtualUsed)
		coll.NetRecv = append(coll.NetRecv, v.NetRecv)
		coll.NetSent = append(coll.NetSent, v.NetSent)
		coll.DiskUsage = append(coll.DiskUsage, v.DiskUsage)
		coll.DiskWrite = append(coll.DiskWrite, v.DiskWrite)
		coll.DiskRead = append(coll.DiskRead, v.DiskRead)
	}

	resp := &GetRealtimeLoadStatResult{}
	resp.MemSwapTotal = m.memSwapTotal
	resp.MemVirtualTotal = m.memVirtualTotal
	resp.DiskTotal = m.diskTotal
	resp.Collection = coll
	return resp, nil
}

func (m *monitor) Collection(ctx context.Context) error {
	for {

		v := &value.DynamicStat{}
		v.Time = time.Now().Unix()
		vv, err := stat.Dynamic()
		if err != nil {
			continue
		}
		v.CpuPercent = formatFloat64(vv.CpuPercent)
		v.Load1 = formatFloat64(vv.Load1)
		v.Load5 = formatFloat64(vv.Load5)
		v.Load15 = formatFloat64(vv.Load15)
		v.MemSwapUsed = vv.MemSwapUsed
		v.MemVirtualUsed = vv.MemVirtualUsed
		v.NetRecv = vv.NetRecv
		v.NetSent = vv.NetSent
		v.DiskUsage = vv.DiskUsage
		v.DiskWrite = vv.DiskWrite
		v.DiskRead = vv.DiskRead

		m.memSwapTotal = vv.MemSwapTotal
		m.memVirtualTotal = vv.MemVirtualTotal
		m.diskTotal = vv.DiskTotal

		m.status = append(m.status, v)
		if len(m.status) > m.maxInterval {
			m.status = m.status[1:]
		}

		time.Sleep(time.Duration(m.interval) * time.Second)
	}
}

func formatFloat64(v float64) float64 {
	vv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v), 64)
	return vv
}
