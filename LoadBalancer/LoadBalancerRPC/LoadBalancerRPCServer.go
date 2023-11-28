package LoadBalancerRPC

import (
	"DistributedFileSystem/LoadBalancer"
)

// LoadBalancerRPCServer Make a LoadBalancerRPCServer to receive RPC request from other places
type LoadBalancerRPCServer struct {
	LoadBalancerServer *LoadBalancer.LoadBalancerServer
}
