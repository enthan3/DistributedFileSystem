package MasterNode

import (
	"DistributedFileSystem/MasterNode/MasterConfiguration"
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/MasterNode/MasterRPC"
	"DistributedFileSystem/MasterNode/MaterLogService"
	"DistributedFileSystem/Metadata"
	"log"
	"net"
	"net/rpc"
	"time"
)

func StartMasterServer() {
	Config, err := MasterConfiguration.LoadMasterConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	Logger, err := MaterLogService.InitLogger(Config.LogPath + "MasterLog")
	if err != nil {
		log.Fatal(err)
	}
	m := MasterRPC.MasterRPCServer{MasterServer: &MasterDefinition.MasterServer{CurrentRPCAddress: Config.CurrentRPCAddress,
		CurrentBackupRPCAddress: Config.MasterBackupRPCAddress, SlaveRPCAddress: Config.SlavesRPCAddress,
		FileMetadataName: make(map[string]*Metadata.FileMetaData), FileMetadataID: make(map[string]*Metadata.FileMetaData),
		ChunkSize: 67108864, StoragePath: Config.StoragePath, ReplicationFactor: 3, LoadBalanceIndex: 0, Logger: Logger,
		IsBackup: Config.IsBackup, LoadBalancerRPCAddress: Config.LoadBalancerRPCAddress}}
	err = rpc.Register(&m)
	if err != nil {
		log.Fatal(err)
	}
	if !Config.IsBackup {
		L, err := net.Listen("tcp", m.MasterServer.CurrentRPCAddress)
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
	} else {
		L, err := net.Listen("tcp", m.MasterServer.CurrentBackupRPCAddress)
		if err != nil {
			log.Fatal(err)
		}
		err = MasterRPC.SendSyncLogToMaster(Logger.LatestID, m.MasterServer)
		if err != nil {
			return
		}
		go func() {
			for {
				err, Status := MasterRPC.SendStatusRequestToMaster(m.MasterServer)
				if err != nil {
					log.Fatal(err)
				}
				if !Status {

				}
				time.Sleep(time.Duration(Config.HeartbeatDuration) * time.Second)
			}
		}()
		go func() {
			for {
				conn, err := L.Accept()
				if err != nil {
					log.Fatal(err)
				}
				go rpc.ServeConn(conn)
			}
		}()
	}

}
