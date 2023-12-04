package SlaveNode

import (
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/SlaveNode/SlaveConfiguration"
	"DistributedFileSystem/SlaveNode/SlaveDefinition"
	"DistributedFileSystem/SlaveNode/SlaveRPC"
	"log"
	"net"
	"net/rpc"
)

func StartSlaveServer() {
	Config, err := SlaveConfiguration.LoadSlaveConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	s := SlaveRPC.SlaveRPCServer{SlaveServer: &SlaveDefinition.SlaveServer{CurrentRPCAddress: Config.CurrentRPCAddress,
		StoragePath: Config.StoragePath, ChunksMetaData: make(map[string]Metadata.ChunkMetaData)}}
	err = rpc.Register(&s)
	if err != nil {
		log.Fatal(err)
	}
	L, err := net.Listen("tcp", s.SlaveServer.CurrentRPCAddress)
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
