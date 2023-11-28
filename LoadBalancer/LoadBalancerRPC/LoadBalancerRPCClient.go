package LoadBalancerRPC

import (
	"DistributedFileSystem/LoadBalancer"
	"DistributedFileSystem/Transmission"
	"errors"
	"net/rpc"
	"syscall"
	"time"
)

// SendStatusRequestToMaster Send Status Request Masters
func SendStatusRequestToMaster(LoadBalancerServer *LoadBalancer.LoadBalancerServer) error {
	var reply Transmission.MasterStatusArg
	for MasterAddress, MasterNode := range LoadBalancerServer.Masters {
		client, err := rpc.Dial("tcp", MasterAddress)
		if err != nil {
			return nil
		}
		//TODO （主节点）实现ReceiveStatusRequest
		err = client.Call("", &struct{}{}, &reply)
		if err != nil {
			if errors.Is(err, rpc.ErrShutdown) || errors.Is(err, syscall.ECONNREFUSED) {
				MasterNode.MasterStatus = &Transmission.MasterStatusArg{MasterAddress: MasterAddress, LastHeartbeatTime: time.Now(), HealthStatus: false}
				return nil
			}
			return err
		}
		MasterNode.MasterStatus = &reply
	}
	return nil
}
