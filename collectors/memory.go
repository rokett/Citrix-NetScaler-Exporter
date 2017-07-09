package collectors

import "github.com/shirou/gopsutil/mem"

// GetMem ...
func GetMem() float64 {
	mem, _ := mem.VirtualMemory()

	return mem.UsedPercent
}
