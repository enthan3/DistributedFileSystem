package MasterDefinition

import "DistributedFileSystem/Metadata"

type MasterServer struct {
	CurrentRPCAddress       string
	CurrentBackupRPCAddress string
	SlaveRPCAddress         []string
	FileMetadata            []Metadata.FileMetaData
	ChunkSize               int64
	StoragePath             string
	ReplicationFactor       int
}
