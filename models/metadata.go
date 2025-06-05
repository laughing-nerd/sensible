package models

type Metadata struct {
	OS               string
	OSVersion        string
	KernelVersion    string
	Architecture     string
	Processor        string
	ProcessorCores   int
	ProcessorThreads int
	MemoryTotal      int64
	MemorySwap       int64
	Hostname         string
	Uptime           uint64 // in seconds
	BootTime         uint64 // in seconds since epoch
	DiskTotal        int64
	DiskFree         int64
}
