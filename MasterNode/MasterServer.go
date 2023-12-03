package MasterNode

import (
	"DistributedFileSystem/MasterNode/MasterConfiguration"
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/MasterNode/MasterRPC"
	"DistributedFileSystem/Metadata"
	"log"
	"net"
	"net/rpc"
)

func StartMasterServer() {
	Config, err := MasterConfiguration.LoadMasterConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	m := MasterRPC.MasterRPCServer{MasterServer: &MasterDefinition.MasterServer{
		CurrentRPCAddress: Config.CurrentRPCAddress, CurrentBackupRPCAddress: Config.MasterBackupRPCAddress,
		SlaveRPCAddress: Config.SlavesRPCAddress, FileMetadata: make([]Metadata.FileMetaData, 0)}}
	err = rpc.Register(&m)
	if err != nil {
		log.Fatal(err)
	}
	L, err := net.Listen("tcp", Config.CurrentRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := L.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeConn(conn)

	}
}
