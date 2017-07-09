package collectors

import "github.com/shirou/gopsutil/cpu"

// GetCPU ...
func GetCPU() float64 {
	//http://<netscaler-ip-address>/nitro/v1/stat/ns
	//cpuusage = cpu usage in percent
	//mgmtcpuusagepcnt = Management CPU utilization percentage.
	//pktcpuusagepcnt = Average CPU utilization percentage for all packet engines excluding management PE.
	//memuseinmb = Main memory currently in use, in megabytes.
	//disk0perusage = Used space in /flash partition of the disk, as a percentage. This is a critical counter. You can configure /flash Used (%) by using the Set snmp alarm DISK-USAGE-HIGH command.
	//disk1perusage = Used space in /var partition of the disk, as a percentage. This is a critical counter. You can configure /var Used (%) by using the Set snmp alarm DISK-USAGE-HIGH command.
	//memusagepcnt = Percentage of memory utilization on NetScaler.
	p, _ := cpu.Percent(0, false)

	return p[0]
}
