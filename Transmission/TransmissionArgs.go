package Transmission

import (
	"github.com/shirou/gopsutil/load"
	"time"
)

type FileArgs struct {
	FileName string
	Size     int64
	Data     []byte
}

type ChunkArgs struct {
	ChunkName string
	Size      int64
	Data      []byte
}

type MasterStatusArg struct {
	MasterAddress     string
	LastHeartbeatTime time.Time
	HealthStatus      bool
	Stat              *load.AvgStat
}
