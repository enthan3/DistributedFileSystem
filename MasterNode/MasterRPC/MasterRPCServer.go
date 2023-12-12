package MasterRPC

import (
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/MasterNode/MasterFileProcessing"
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/Transmission"
	"DistributedFileSystem/Utils"
	"github.com/shirou/gopsutil/load"
	"time"
)

type MasterRPCServer struct {
	MasterServer *MasterDefinition.MasterServer
}

func (m *MasterRPCServer) ReceiveStatusRequestFromLoadBalancer(none *struct{}, reply *Transmission.MasterStatusArg) error {
	AvgStat, err := load.Avg()
	if err != nil {
		return err
	}
	MasterStatus := Transmission.MasterStatusArg{MasterAddress: m.MasterServer.CurrentRPCAddress, LastHeartbeatTime: time.Now(), HealthStatus: true, Stat: AvgStat}
	*reply = MasterStatus
	return nil
}

func (m *MasterRPCServer) ReceiveFileFromFrontendService(FileArg *Transmission.FileArgs, reply *bool) error {
	var FileMetadata Metadata.FileMetaData
	FileUUID, _ := Utils.GenerateUniqueID()
	FileChunks, err := MasterFileProcessing.ChunkFile(FileUUID, FileArg, m.MasterServer)
	if err != nil {
		return err
	}
	FileChunkReplicates, err := MasterFileProcessing.ReplicationChunks(FileChunks, m.MasterServer)
	err = SendChunksToSlave(&FileChunkReplicates, &FileMetadata, m.MasterServer)
	if err != nil {
		return err
	}
	FileMetadata.FileID = FileUUID
	FileMetadata.FileName = FileArg.FileName
	FileMetadata.Size = FileArg.Size
	FileMetadata.CreateTime = time.Now()
	m.MasterServer.FileMetadataName[FileArg.FileName] = &FileMetadata
	m.MasterServer.FileMetadataID[FileArg.FileName] = &FileMetadata
	*reply = true
	return nil
}

func (m *MasterRPCServer) ReceiveDeleteFromFrontendService(FileName string, reply *string) error {

	err := SendDeleteToSlave(FileName, m.MasterServer)
	if err != nil && err.Error() == "File does not exist!" {
		*reply = "File does not exist!"
		return nil
	} else {
		*reply = "Delete send to Slave Master error"
	}
	delete(m.MasterServer.FileMetadataName, FileName)
	delete(m.MasterServer.FileMetadataID, FileName)
	*reply = "Success!"
	return nil
}

func (m *MasterRPCServer) ReceiveSearchFromFrontendService(Filename string, reply *Metadata.FileMetaData) error {
	FileMetadata, _ := m.MasterServer.FileMetadataName[Filename]
	*reply = *FileMetadata
	return nil
}
