//go:build linux

package runtime

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func (l *Local) Snapshot() (Snapshot, error) {
	cpu, err := cpuPercent()
	if err != nil {
		return Snapshot{}, err
	}
	memory, err := memoryPercent()
	if err != nil {
		return Snapshot{}, err
	}
	disk, err := diskPercent(l.rootPath)
	if err != nil {
		return Snapshot{}, err
	}
	uptime, err := uptimeSeconds()
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		CPUPercent:     cpu,
		MemoryPercent:  memory,
		DiskPercent:    disk,
		ActiveSessions: 0,
		UptimeSeconds:  uptime,
	}, nil
}

type cpuStat struct {
	total uint64
	idle  uint64
}

func cpuPercent() (float64, error) {
	first, err := readCPUStat()
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	second, err := readCPUStat()
	if err != nil {
		return 0, err
	}
	totalDelta := second.total - first.total
	idleDelta := second.idle - first.idle
	if totalDelta == 0 {
		return 0, nil
	}
	value := float64(totalDelta-idleDelta) / float64(totalDelta) * 100
	if value < 0 {
		return 0, nil
	}
	if value > 100 {
		return 100, nil
	}
	return value, nil
}

func readCPUStat() (cpuStat, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return cpuStat{}, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] != "cpu" {
			continue
		}
		var total uint64
		for _, field := range fields[1:] {
			n, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return cpuStat{}, err
			}
			total += n
		}
		idle, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return cpuStat{}, err
		}
		if len(fields) > 5 {
			iowait, err := strconv.ParseUint(fields[5], 10, 64)
			if err != nil {
				return cpuStat{}, err
			}
			idle += iowait
		}
		return cpuStat{total: total, idle: idle}, nil
	}
	return cpuStat{}, errors.New("aggregate CPU counter not found")
}

func memoryPercent() (float64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	var total, available uint64
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total = value
		case "MemAvailable:":
			available = value
		}
	}
	if total == 0 {
		return 0, errors.New("memory total not found")
	}
	used := total - minUint64(total, available)
	return float64(used) / float64(total) * 100, nil
}

func diskPercent(path string) (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}
	if stat.Blocks == 0 {
		return 0, errors.New("filesystem has zero blocks")
	}
	used := stat.Blocks - stat.Bfree
	return float64(used) / float64(stat.Blocks) * 100, nil
}

func uptimeSeconds() (int64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, errors.New("uptime not found")
	}
	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}
	return int64(value), nil
}

func minUint64(a, b uint64) uint64 {
	if b < a {
		return b
	}
	return a
}
