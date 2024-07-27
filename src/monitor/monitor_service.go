package monitor

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"dxkite.cn/meownest/pkg/stat"
)

type Monitor interface {
	Collection(ctx context.Context) error
	ListDynamicStat(ctx context.Context, param *ListDynamicStatRequest) (*DynamicStatResult, error)
}

type MonitorConfig struct {
	// 统计间隔
	Interval int
	// 实时数据保留时长
	MaxInterval int
	// 聚合间隔
	RollInterval int
}

type monitor struct {
	interval        int
	maxInterval     int
	rollInterval    int
	memSwapTotal    uint64
	memVirtualTotal uint64
	diskTotal       uint64
	status          []*DynamicStat
	r               DynamicStatRepository
	roll            []*DynamicStat
	mtx             *sync.Mutex
}

func NewMonitor(cfg *MonitorConfig, r DynamicStatRepository) Monitor {
	m := &monitor{r: r}
	m.interval = cfg.Interval
	m.maxInterval = cfg.MaxInterval
	m.rollInterval = cfg.RollInterval
	m.mtx = &sync.Mutex{}
	return m
}

type DynamicStatResult struct {
	Collection      *DynamicStatCollection `json:"collection"`
	MemSwapTotal    uint64                 `json:"mem_swap_total"`
	MemVirtualTotal uint64                 `json:"mem_virtual_total"`
	DiskTotal       uint64                 `json:"disk_total"`
}

type ListDynamicStatRequest struct {
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
}

func (m monitor) ListDynamicStat(ctx context.Context, param *ListDynamicStatRequest) (*DynamicStatResult, error) {
	var startTime, endTime uint64

	if param.StartTime != "" {
		v, err := time.Parse(time.RFC3339, param.StartTime)
		if err != nil {
			return nil, err
		}
		startTime = uint64(v.Unix())
	} else {
		startTime = uint64(time.Now().Add(time.Duration(-m.maxInterval) * time.Second).Unix())
	}

	if param.EndTime != "" {
		v, err := time.Parse(time.RFC3339, param.EndTime)
		if err != nil {
			return nil, err
		}
		endTime = uint64(v.Unix())
	} else {
		endTime = uint64(time.Now().Unix())
	}

	realTimeStart := uint64(time.Now().Unix())
	if len(m.status) > 0 {
		realTimeStart = m.status[0].Time
	}

	// 取实时数据
	output := []*DynamicStat{}
	for _, v := range m.status {
		if v.Time < startTime {
			continue
		}
		if v.Time > endTime {
			continue
		}
		output = append(output, v)
	}

	// 取历史数据 -> 实时数据

	if startTime < realTimeStart {
		entities, err := m.r.List(ctx, &ListDynamicStatParam{
			StartTime: startTime,
			EndTime:   realTimeStart,
		})
		if err != nil {
			return nil, err
		}

		output = append(entities, output...)
	}

	resp := &DynamicStatResult{}
	resp.Collection = NewDynamicStatCollection(output, m.memSwapTotal, m.memVirtualTotal, m.diskTotal)
	resp.MemSwapTotal = m.memSwapTotal
	resp.MemVirtualTotal = m.memVirtualTotal
	resp.DiskTotal = m.diskTotal
	return resp, nil
}

func (s *monitor) Collection(ctx context.Context) error {
	for {

		v := &DynamicStat{}
		v.Time = uint64(time.Now().Unix())
		vv, err := stat.Dynamic()
		if err != nil {
			continue
		}
		v.CpuPercent = vv.CpuPercent
		v.Load1 = vv.Load1
		v.Load5 = vv.Load5
		v.Load15 = vv.Load15
		v.MemSwapUsed = vv.MemSwapUsed
		v.MemVirtualUsed = vv.MemVirtualUsed
		v.NetRecv = vv.NetRecv
		v.NetSent = vv.NetSent
		v.DiskUsage = vv.DiskUsage
		v.DiskWrite = vv.DiskWrite
		v.DiskRead = vv.DiskRead

		s.collect(ctx, v)

		s.memSwapTotal = vv.MemSwapTotal
		s.memVirtualTotal = vv.MemVirtualTotal
		s.diskTotal = vv.DiskTotal

		time.Sleep(time.Duration(s.interval) * time.Second)
	}
}

func (s *monitor) collect(ctx context.Context, ent *DynamicStat) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.status = append(s.status, ent)
	s.roll = append(s.roll, ent)

	if len(s.status) > 0 && ent.Time-s.status[0].Time >= uint64(s.maxInterval) {
		s.status = s.status[1:]
	}

	if len(s.status) > 0 && ent.Time-s.roll[0].Time >= uint64(s.rollInterval) {
		s.rollCollect(ctx, s.roll)
		s.roll = s.roll[0:0]
	}
}

func (s *monitor) rollCollect(ctx context.Context, entities []*DynamicStat) error {
	avg := &DynamicStat{}
	n := len(entities)
	end := entities[n-1]

	for _, v := range entities {
		avg.CpuPercent += v.CpuPercent
		avg.Load1 += v.Load1
		avg.Load5 += v.Load5
		avg.Load15 += v.Load15
	}

	avg.Time = end.Time
	avg.CpuPercent = formatFloat64(avg.CpuPercent / float64(n))
	avg.Load1 = formatFloat64(avg.Load1 / float64(n))
	avg.Load5 = formatFloat64(avg.Load5 / float64(n))
	avg.Load15 = formatFloat64(avg.Load15 / float64(n))
	avg.MemSwapUsed = end.MemSwapUsed
	avg.MemVirtualUsed = end.MemVirtualUsed
	avg.NetRecv = end.NetRecv
	avg.NetSent = end.NetSent
	avg.DiskUsage = end.DiskUsage
	avg.DiskWrite = end.DiskWrite
	avg.DiskRead = end.DiskRead
	_, err := s.r.Create(ctx, avg)
	return err
}

func formatFloat64(v float64) float64 {
	vv, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", v), 64)
	return vv
}
