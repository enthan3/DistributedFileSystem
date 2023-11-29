package LoadBalancerRPC

import (
	"DistributedFileSystem/LoadBalancer"
	"DistributedFileSystem/Transmission"
	"errors"
)

// LoadBalancerRPCServer Make a LoadBalancerRPCServer to receive RPC request from other places
type LoadBalancerRPCServer struct {
	LoadBalancerServer *LoadBalancer.LoadBalancerServer
}

// ReceivePromotionFromBackup Receive Promotion Request from Master Backup Server to handle Master Node Request
func (l *LoadBalancerRPCServer) ReceivePromotionFromBackup(MasterBackupAddress string, reply *bool) error {
	for MasterAddress, MasterBackup := range l.LoadBalancerServer.MasterBackups {
		if MasterBackup == MasterBackupAddress {
			delete(l.LoadBalancerServer.MasterBackups, MasterAddress)
			l.LoadBalancerServer.MasterBackups[MasterBackup] = MasterAddress
			delete(l.LoadBalancerServer.Masters, MasterAddress)
			l.LoadBalancerServer.Masters[MasterBackup] = LoadBalancer.MasterNode{MasterAddress: MasterBackup, MasterStatus: new(Transmission.MasterStatusArg)}
			*reply = true
			return nil
		}
	}
	*reply = false
	return errors.New("Master/MasterBackup not exist!")
}
