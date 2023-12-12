package MasterDefinition

import "DistributedFileSystem/Metadata"

type MasterServer struct {
	CurrentRPCAddress       string
	CurrentBackupRPCAddress string
	SlaveRPCAddress         []string
	FileMetadataName        map[string]*Metadata.FileMetaData
	FileMetadataID          map[string]*Metadata.FileMetaData
	ChunkSize               int64
	StoragePath             string
	ReplicationFactor       int
	LoadBalanceIndex        int
	IsBackup                bool
}
