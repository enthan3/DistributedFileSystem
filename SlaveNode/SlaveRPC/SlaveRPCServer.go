package SlaveRPC

import (
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/SlaveNode/SlaveDefinition"
	"DistributedFileSystem/Transmission"
	"os"
	"time"
)

type SlaveRPCServer struct {
	SlaveServer *SlaveDefinition.SlaveServer
}

func (s *SlaveRPCServer) ReceiveChunkFromMaster(Chunk *Transmission.ChunkArg, reply *bool) error {
	var ChunkMetadata Metadata.ChunkMetaData
	ChunkMetadata.ChunkNode = s.SlaveServer.CurrentRPCAddress
	ChunkMetadata.ChunkName = Chunk.ChunkName
	ChunkMetadata.Size = Chunk.Size
	ChunkMetadata.CreateTime = time.Now()
	err := os.WriteFile(s.SlaveServer.StoragePath+Chunk.ChunkName, Chunk.Data, 0666)
	if err != nil {
		return err
	}
	s.SlaveServer.ChunksMetaData[Chunk.ChunkName] = ChunkMetadata
	*reply = true
	return nil
}

func (s *SlaveRPCServer) ReceiveDeleteFromMaster(ChunkName string, reply *bool) error {
	err := os.Remove(s.SlaveServer.StoragePath + ChunkName)
	if err != nil {
		return err
	}
	delete(s.SlaveServer.ChunksMetaData, ChunkName)
	*reply = true
	return nil
}

func (s *SlaveRPCServer) ReceiveDownloadFromFrontendService(ChunkName string, Chunk *Transmission.ChunkArg) error {
	Data, err := os.ReadFile(s.SlaveServer.StoragePath + ChunkName)
	if err != nil {
		return err
	}
	ChunkMetadata := s.SlaveServer.ChunksMetaData[ChunkName]
	Chunk.Size = ChunkMetadata.Size
	Chunk.Data = Data
	Chunk.ChunkName = ChunkMetadata.ChunkName
	return nil
}
