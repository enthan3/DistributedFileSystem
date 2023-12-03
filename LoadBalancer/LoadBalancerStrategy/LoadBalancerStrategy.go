package LoadBalancerStrategy

import (
	"DistributedFileSystem/LoadBalancer/LoadBalancerDefinition"
	"math"
)

// GetNextMasterURL Use the lowest CPU and memory usage to determine the next request should be hand into which Master Server
func GetNextMasterURL(l *LoadBalancerDefinition.LoadBalancerServer) string {
	LowestLoad := math.MaxFloat64
	LowestAddress := ""
	for MasterAddress, MasterStatus := range l.MasterStatusMap {
		if MasterStatus.HealthStatus == true && MasterStatus.Stat.Load1 < LowestLoad {
			LowestLoad = MasterStatus.Stat.Load1
			LowestAddress = MasterAddress
		}
	}
	return LowestAddress
}
