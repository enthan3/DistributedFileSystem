package SlaveDefinition

import "DistributedFileSystem/Metadata"

type SlaveServer struct {
	CurrentRPCAddress string
	StoragePath       string
	ChunksMetaData    map[string]Metadata.ChunkMetaData
}
