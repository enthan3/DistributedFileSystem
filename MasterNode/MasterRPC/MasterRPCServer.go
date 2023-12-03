package MasterRPC

import (
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/Transmission"
	"github.com/shirou/gopsutil/load"
	"time"
)

type MasterRPCServer struct {
	MasterServer *MasterDefinition.MasterServer
}

func (m *MasterRPCServer) ReceiveStatusRequestFromLoadBalancer(none *struct{}, reply *Transmission.MasterStatusArg) error {
	AvgStat, err := load.Avg()
	if err != nil {
		return err
	}
	MasterStatus := Transmission.MasterStatusArg{MasterAddress: m.MasterServer.CurrentRPCAddress, LastHeartbeatTime: time.Now(), HealthStatus: true, Stat: AvgStat}
	*reply = MasterStatus
	return nil
}

func (m *MasterRPCServer) ReceiveFileFromFrontendService(FileArg *Transmission.FileArgs, reply *bool) error {
	//TODO
	return nil
}

func (m *MasterRPCServer) ReceiveDeleteFromFrontendService(FileName string, reply *bool) error {
	//TODO
	return nil
}

func (m *MasterRPCServer) ReceiveSearchFromFrontendService() error {
	//TODO
	return nil
}
