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

// SendMasterToFrontendService Send Frontend Service with Master Node address with the lowest load to processing request
func SendMasterToFrontendService(l *LoadBalancerDefinition.LoadBalancerServer) error {
	var reply bool
	MasterAddress := LoadBalancerStrategy.GetNextMasterURL(l)
	client, err := rpc.Dial("tcp", MasterAddress)
	if err != nil {
		return err
	}
	//TODO 前端服务层实现接收主节点,实现如果reply不为正确的选项
	err = client.Call("", &MasterAddress, &reply)
	if err != nil {
		return err
	}
	return nil
}
