package stat

import (
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// 动态数据
type DynamicStat struct {
	CpuPercent      float64 `json:"cpu_percent"`
	Load1           float64 `json:"load_1"`
	Load5           float64 `json:"load_5"`
	Load15          float64 `json:"load_15"`
	MemSwapUsed     uint64  `json:"mem_swap_used"`
	MemSwapTotal    uint64  `json:"mem_swap_total"`
	MemVirtualUsed  uint64  `json:"mem_virtual_used"`
	MemVirtualTotal uint64  `json:"mem_virtual_total"`
	NetRecv         uint64  `json:"net_recv"`
	NetSent         uint64  `json:"net_send"`
	DiskUsage       uint64  `json:"disk_usage"`
	DiskTotal       uint64  `json:"disk_total"`
	DiskWrite       uint64  `json:"disk_write"`
	DiskRead        uint64  `json:"disk_read"`
}

// 统计动态数据
func Dynamic() (*DynamicStat, error) {
	stat := &DynamicStat{}

	if v, err := cpu.Percent(0, false); err != nil {
		return nil, err
	} else {
		stat.CpuPercent = v[0]
	}

	if v, err := load.Avg(); err != nil {
		return nil, err
	} else {
		stat.Load1 = v.Load1
		stat.Load5 = v.Load5
		stat.Load15 = v.Load15
	}

	if v, err := mem.SwapMemory(); err != nil {
		return nil, err
	} else {
		stat.MemSwapUsed = v.Used
		stat.MemSwapTotal = v.Total
	}

	if v, err := mem.VirtualMemory(); err != nil {
		return nil, err
	} else {
		stat.MemVirtualUsed = v.Used
		stat.MemVirtualTotal = v.Total
	}

	if send, recv, err := getNetStatus(); err != nil {
		return nil, err
	} else {
		stat.NetRecv = recv
		stat.NetSent = send
	}

	if total, used, err := getDiskStatus(); err != nil {
		return nil, err
	} else {
		stat.DiskUsage = used
		stat.DiskTotal = total
	}

	if read, write, err := getDiskIoStatus(); err != nil {
		return nil, err
	} else {
		stat.DiskRead = read
		stat.DiskWrite = write
	}

	return stat, nil
}

func getNetStatus() (send, recv uint64, err error) {
	if v, e1 := net.IOCounters(true); err != nil {
		err = e1
		return
	} else {
		for _, iter := range v {
			if !isSkipInterface(iter.Name) {
				send += iter.BytesSent
				recv += iter.BytesRecv
			}
		}
	}
	return
}

func getDiskStatus() (total, used uint64, err error) {
	parts, err := disk.Partitions(false)
	for _, part := range parts {
		if !isKnownFs(part.Fstype) {
			continue
		}
		if usage, err := disk.Usage(part.Mountpoint); err == nil {
			total += usage.Total
			used += usage.Used
		}
	}
	return
}

func getDiskIoStatus() (read, write uint64, err error) {
	data, _ := disk.IOCounters()
	for _, v := range data {
		read += v.ReadBytes
		write += v.WriteBytes
	}
	return
}

type SystemStat struct {
	Hostname             string `json:"hostname"`
	BootTime             uint64 `json:"boot_time"`
	CpuCount             int    `json:"cpu_count"`
	OS                   string `json:"os"`
	Platform             string `json:"platform"`
	PlatformFamily       string `json:"platform_family"`
	PlatformVersion      string `json:"platform_version"`
	KernelVersion        string `json:"kernel_version"`
	KernelArch           string `json:"kernel_arch"`
	VirtualizationSystem string `json:"virtualization_system"`
	VirtualizationRole   string `json:"virtualization_role"`
}

// 获取系统信息
func System() (*SystemStat, error) {
	stat := &SystemStat{}
	if v, err := host.Info(); err != nil {
		return nil, err
	} else {
		stat.Hostname = v.Hostname
		stat.BootTime = v.BootTime
		stat.OS = v.OS
		stat.Platform = v.Platform
		stat.PlatformFamily = v.PlatformFamily
		stat.PlatformVersion = v.PlatformVersion
		stat.KernelVersion = v.KernelVersion
		stat.VirtualizationSystem = v.VirtualizationSystem
		stat.VirtualizationRole = v.VirtualizationRole
	}

	if v, err := cpu.Counts(true); err != nil {
		return nil, err
	} else {
		stat.CpuCount = v
	}

	return stat, nil
}

func isSkipInterface(name string) bool {
	name = strings.ToLower(name)
	for _, v := range skipInterfaceNames {
		if strings.Contains(name, v) {
			return false
		}
	}
	return false
}

func isKnownFs(name string) bool {
	name = strings.ToLower(name)
	for _, v := range knownFsNames {
		if name == v {
			return true
		}
	}
	return false
}

var skipInterfaceNames = []string{"lo", "tun", "kube", "docker", "vmbr", "br-", "vnet", "veth"}
var knownFsNames = []string{"ext4", "ext3", "ext2", "reiserfs", "jfs", "btrfs", "fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs", "apfs"}
