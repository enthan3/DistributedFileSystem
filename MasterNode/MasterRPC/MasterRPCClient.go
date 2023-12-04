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
	FileMetadata, exist := m.FileMetadata[Filename]
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
