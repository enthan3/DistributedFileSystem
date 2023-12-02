package LoadBalancerRPC

import (
	"DistributedFileSystem/LoadBalancer/LoadBalancerDefinition"
	"DistributedFileSystem/Transmission"
	"errors"
)

// LoadBalancerRPCServer Make a LoadBalancerRPCServer to receive RPC request from other places
type LoadBalancerRPCServer struct {
	LoadBalancerServer *LoadBalancerDefinition.LoadBalancerServer
}

// ReceivePromotionFromBackup Receive Promotion Request from Master Backup Server receiving
func (l *LoadBalancerRPCServer) ReceivePromotionFromBackup(MasterBackupAddress *string, reply *bool) error {
	for MasterAddress, MasterBackup := range l.LoadBalancerServer.MasterBackupsMap {
		if MasterBackup == *MasterBackupAddress {
			delete(l.LoadBalancerServer.MasterBackupsMap, MasterAddress)
			l.LoadBalancerServer.MasterBackupsMap[MasterBackup] = MasterAddress
			delete(l.LoadBalancerServer.MasterStatusMap, MasterAddress)
			l.LoadBalancerServer.MasterStatusMap[MasterBackup] = &Transmission.MasterStatusArg{}
			err := SendMastersToFrontendService(l.LoadBalancerServer)
			if err != nil {
				return err
			}
			*reply = true
			return nil
		}
	}
	*reply = false
	return errors.New("Master/MasterBackup not exist!")
}
