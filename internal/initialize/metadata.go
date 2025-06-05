package initialize

import (
	"sensible/models"

	"github.com/shirou/gopsutil/v4/host"
)

func FetchMetadata() {
	var metadata models.Metadata

	hostInfo, err := host.Info()
	if err != nil {
		return
	}

	if hostInfo != nil {
		metadata.OS = hostInfo.OS
		metadata.OSVersion = hostInfo.PlatformVersion
		metadata.KernelVersion = hostInfo.KernelVersion
		metadata.Architecture = hostInfo.KernelArch
		metadata.Hostname = hostInfo.Hostname
		metadata.Uptime = hostInfo.Uptime
		metadata.BootTime = hostInfo.BootTime
	}

}
