package LoadBalancerStrategy

import "DistributedFileSystem/LoadBalancer"

// GetNextMasterURL Use the lowest CPU and memory usage to determine the next request should be hand into which Master Server
func GetNextMasterURL(l *LoadBalancer.LoadBalancerServer) string {

	return ""
}
