package MasterDefinition

import "DistributedFileSystem/Metadata"

type MasterServer struct {
	CurrentRPCAddress       string
	CurrentBackupRPCAddress string
	SlaveRPCAddress         []string
	FileMetadata            []Metadata.FileMetaData
}

type MasterBackupServer struct {
	CurrentRPCAddress string
	MasterRPCAddress  string
	SlaveRPCAddress   []string
	FileMetadata      []Metadata.FileMetaData
}
