package MasterRPC

import (
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/MasterNode/MasterLoadBalance"
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/Transmission"
	"errors"
	"net/rpc"
	"strings"
	"time"
)

//Master Send Files Functions

func SendChunksToSlave(chunkArgs *[]Transmission.ChunkArg, FileMetaData *Metadata.FileMetaData, m *MasterDefinition.MasterServer) error {
	var reply bool
	for _, Chunk := range *chunkArgs {
		Slave := MasterLoadBalance.ChooseSlaveRoundRobin(m)
		client, err := rpc.Dial("tcp", Slave)
		if err != nil {
			return err
		}
		err = client.Call("SlaveRPCServer.ReceiveChunkFromMaster", &Chunk, &reply)
		if err != nil {
			return err
		}
		if !reply {
			return errors.New("File Send Master error")
		}
		ChunkMetaData := Metadata.ChunkMetaData{ChunkName: Chunk.ChunkName, Size: Chunk.Size, CreateTime: time.Now(), ChunkNode: Slave}
		if strings.Contains(Chunk.ChunkName, "main") {
			FileMetaData.ChunkMain = append(FileMetaData.ChunkMain, ChunkMetaData)
		} else if strings.Contains(Chunk.ChunkName, "replicate") {
			FileMetaData.ChunkReplicate = append(FileMetaData.ChunkReplicate, ChunkMetaData)
		}
	}
	return nil
}

func SendDeleteToSlave(Filename string, m *MasterDefinition.MasterServer) error {
	var reply bool
	FileMetadata, exist := m.FileMetadataName[Filename]
	if exist {
		for _, ChunkMetadata := range FileMetadata.ChunkMain {
			client, err := rpc.Dial("tcp", ChunkMetadata.ChunkNode)
			if err != nil {
				return err
			}
			err = client.Call("SlaveRPCServer.ReceiveDeleteFromMaster", ChunkMetadata.ChunkName, &reply)
			if err != nil {
				return err
			}
			if !reply {
				return errors.New("File Delete Master error")
			}

		}
		for _, ChunkReplicateMetadata := range FileMetadata.ChunkReplicate {
			client, err := rpc.Dial("tcp", ChunkReplicateMetadata.ChunkNode)
			if err != nil {
				return err
			}
			err = client.Call("SlaveRPCServer.ReceiveDeleteFromMaster", ChunkReplicateMetadata.ChunkName, &reply)
			if err != nil {
				return err
			}
			if !reply {
				return errors.New("File Delete Master error")
			}

		}
		return nil
	} else {
		return errors.New("File does not exist!")
	}
}

// Master Sync Log Functions
func SendSyncLogToMaster(LatestID int64, m *MasterDefinition.MasterServer) error {
	var reply bool
	client, err := rpc.Dial("tcp", m.CurrentRPCAddress)
	if err != nil {
		return err
	}
	err = client.Call("MasterRPCServer.ReceiveSyncLogFromMasterBackup", LatestID, &reply)
	if !reply {
		return errors.New("Send sync log MasterBackup error")
	}
	return nil
}

func SendStatusRequestToMaster(m *MasterDefinition.MasterServer) (error, bool) {
	var reply Transmission.MasterStatusArg
	client, err := rpc.Dial("tcp", m.CurrentRPCAddress)
	if err != nil {
		return nil, false
	}
	err = client.Call("MasterRPCServer.ReceiveStatusRequestFromLoadBalancer", &struct{}{}, &reply)
	if !reply.HealthStatus {
		return nil, false
	}
	return nil, true
}

func SendPromotionRequestToLoadBalancer(m *MasterDefinition.MasterServer) error {
	var reply bool
	client, err := rpc.Dial("tcp", m.LoadBalancerRPCAddress)
	if err != nil {
		return err
	}
	err = client.Call("LoadBalancerRPCServer.ReceivePromotionFromBackup", m.CurrentBackupRPCAddress, &reply)
	if err != nil {
		return err
	}
	if !reply {
		return errors.New("Send promotino request to Master error")
	}
	return nil
}

// Master To MasterBackups Updates FileMetadata function
func SendFileMetadataToMasterBackup(FileMetadata *Metadata.FileMetaData, m *MasterDefinition.MasterServer) {
	var reply bool
	client, err := rpc.Dial("tcp", m.CurrentBackupRPCAddress)
	if err != nil {
		return
	}
	//TODO
	err = client.Call("MasterRPCServer.ReceiveFileMetadataFromMaster", FileMetadata, &reply)
	if err != nil {
		return
	}
	return
}

func SendSearchToMasterBackup(FileUUID string, m *MasterDefinition.MasterServer) {
	var reply bool
	client, err := rpc.Dial("tcp", m.CurrentBackupRPCAddress)
	if err != nil {
		return
	}
	//TODO
	err = client.Call("MasterRPCServer.ReceiveSearchFromMaster", FileUUID, &reply)
	if err != nil {
		return
	}
	return
}

func SendDeleteToMasterBackup(FileMetadata *Metadata.FileMetaData, m *MasterDefinition.MasterServer) {
	var reply bool
	client, err := rpc.Dial("tcp", m.CurrentBackupRPCAddress)
	if err != nil {
		return
	}
	//TODO
	err = client.Call("MasterRPCServer.ReceiveDeleteFromMaster", FileMetadata, &reply)
	if err != nil {
		return
	}
	return
}
