package FrontendServiceRPC

import "DistributedFileSystem/FrontendService/FrontendServiceDefinition"

type FrontendServiceRPCServer struct {
	FrontendServiceServer *FrontendServiceDefinition.FrontendServiceServer
}

func (f *FrontendServiceRPCServer) ReceiveLowestUsageMasterFromLoadBalancer(MasterAddress *string, reply *bool) error {
	f.FrontendServiceServer.LowestMaster = *MasterAddress
	*reply = true
	return nil
}

func (f *FrontendServiceRPCServer) ReceiveMastersFromLoadBalancer(MastersAddress *[]string, reply *bool) error {
	f.FrontendServiceServer.Masters = *MastersAddress
	*reply = true
	return nil
}
