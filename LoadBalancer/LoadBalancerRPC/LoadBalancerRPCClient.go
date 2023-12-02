package LoadBalancerRPC

import (
	"DistributedFileSystem/LoadBalancer/LoadBalancerDefinition"
	"DistributedFileSystem/LoadBalancer/LoadBalancerStrategy"
	"DistributedFileSystem/Transmission"
	"errors"
	"net/rpc"
	"time"
)

// SendStatusRequestToMaster Send Status Request Masters
func SendStatusRequestToMaster(l *LoadBalancerDefinition.LoadBalancerServer) error {
	var reply Transmission.MasterStatusArg
	for MasterAddress, _ := range l.MasterStatusMap {
		client, err := rpc.Dial("tcp", MasterAddress)
		if err != nil {
			return err
		}
		//TODO 主节点实现状态更新检查,可能实现reply不为正确的选项
		err = client.Call("", &struct{}{}, &reply)
		if err != nil {
			if errors.Is(err, rpc.ErrShutdown) {
				l.MasterStatusMap[MasterAddress] = &Transmission.MasterStatusArg{MasterAddress: MasterAddress, LastHeartbeatTime: time.Now(), HealthStatus: false}
				return nil
			}
			return err
		}
		l.MasterStatusMap[MasterAddress] = &reply
	}
	return nil
}

// SendLowestUsageMasterToFrontendService Send Frontend Service with Master Node address with the lowest load to processing request
func SendLowestUsageMasterToFrontendService(l *LoadBalancerDefinition.LoadBalancerServer) error {
	var reply bool
	MasterAddress := LoadBalancerStrategy.GetNextMasterURL(l)
	client, err := rpc.Dial("tcp", MasterAddress)
	if err != nil {
		return err
	}
	err = client.Call("FrontendServiceRPCServer.ReceiveLowestUsageMasterFromLoadBalancer", &MasterAddress, &reply)
	if err != nil {
		return err
	}
	if reply != true {
		return errors.New("Receive lowest usage master address FrontendService error")
	}
	return nil
}

// SendMastersToFrontendService Send Frontend Service with all the master node address
func SendMastersToFrontendService(l *LoadBalancerDefinition.LoadBalancerServer) error {
	var reply bool
	var Masters []string

	for MasterAddress, _ := range l.MasterBackupsMap {
		Masters = append(Masters, MasterAddress)
	}
	client, err := rpc.Dial("tcp", l.ServiceRPC)
	if err != nil {
		return err
	}
	err = client.Call("FrontendServiceRPCServer.ReceiveMastersFromLoadBalancer", &Masters, &reply)
	if err != nil {
		return err
	}
	if reply != true {
		return errors.New("Receive masters address FrontendService error")
	}
	return nil
}
