package LoadBalancer

import (
	"DistributedFileSystem/Transmission"
)

type MasterNode struct {
	MasterAddress string
	MasterStatus  *Transmission.MasterStatusArg
}

type LoadBalancerServer struct {
	//Map structure Master Address to MasterNode
	Masters map[string]MasterNode
	//Map structure Master Address to Master Backup
	MasterBackups map[string]string
	//Service Address
	Service string
}

func StartLoadBalancerServer() {

}
