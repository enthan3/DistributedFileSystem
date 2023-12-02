package LoadBalancerDefinition

import "DistributedFileSystem/Transmission"

type LoadBalancerServer struct {
	//Map structure Master Address to MasterNode
	MasterStatusMap map[string]*Transmission.MasterStatusArg
	//Map structure Master Address to Master Backup
	MasterBackupsMap map[string]string
	//Service Address
	ServiceHTTP    string
	ServiceRPC     string
	CurrentAddress string
}
