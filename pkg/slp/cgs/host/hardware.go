package host

import (
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/shirou/gopsutil/v4/disk"
)

type Memory struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Percent   float64 `json:"percent"`
}

type Disk struct {
	Total   uint64  `json:"total"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

type CPU struct {
	Percent float64 `json:"percent"`

	User    float64 `json:"user"`
	Nice    float64 `json:"nice"`
	System  float64 `json:"system"`
	Idle    float64 `json:"idle"`
	IoWait  float64 `json:"io_wait"`
	Irq     float64 `json:"irq"`
	SoftIrq float64 `json:"soft_irq"`

	CacheSize uint64   `json:"cache_size"`
	ModelName string   `json:"model_name"`
	Family    string   `json:"family"`
	Model     string   `json:"model"`
	Stepping  string   `json:"stepping"`
	Flags     []string `json:"flags"`
	MHz       uint64   `json:"mhz"`
}

type Hardware struct {
	Memory   Memory `json:"memory"`
	MainDisk Disk   `json:"main_disk"`
	CPU      []CPU  `json:"cpu"`
}

func GetHardwareProfile() (*Hardware, error) {

	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	diskUsage, err := disk.Usage(homeDir)
	if err != nil {
		return nil, err
	}

	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	cpuTimeStats, err := cpu.Times(true)
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpus := []CPU{}
	for idx, cpu := range cpuTimeStats {

		cpu := CPU{
			User:    cpu.User,
			Nice:    cpu.Nice,
			System:  cpu.System,
			Idle:    cpu.Idle,
			IoWait:  cpu.Iowait,
			Irq:     cpu.Irq,
			SoftIrq: cpu.Softirq,
		}

		if len(cpuUsage) > idx {
			cpu.Percent = cpuUsage[idx]
		} else if len(cpuUsage) > 0 {
			cpu.Percent = cpuUsage[0]
		}

		if len(cpuInfo) > idx {
			cpu.CacheSize = uint64(cpuInfo[idx].CacheSize)
			cpu.ModelName = cpuInfo[idx].ModelName
			cpu.Family = cpuInfo[idx].Family
			cpu.Model = cpuInfo[idx].Model
			cpu.Stepping = strconv.Itoa(int(cpuInfo[idx].Stepping))
			cpu.MHz = uint64(cpuInfo[idx].Mhz)
		} else if len(cpuInfo) > 0 {
			cpu.CacheSize = uint64(cpuInfo[0].CacheSize)
			cpu.ModelName = cpuInfo[0].ModelName
			cpu.Family = cpuInfo[0].Family
			cpu.Model = cpuInfo[0].Model
			cpu.Stepping = strconv.Itoa(int(cpuInfo[0].Stepping))
			cpu.MHz = uint64(cpuInfo[0].Mhz)
		}

		cpus = append(cpus, cpu)
	}

	return &Hardware{
		Memory: Memory{
			Total:     memory.Total,
			Available: memory.Available,
			Used:      memory.Used,
			Percent:   memory.UsedPercent,
		},
		MainDisk: Disk{
			Total:   diskUsage.Total,
			Used:    diskUsage.Used,
			Percent: diskUsage.UsedPercent,
		},
		CPU: cpus,
	}, nil
}
