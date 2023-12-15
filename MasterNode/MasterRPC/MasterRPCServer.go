package MasterRPC

import (
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/MasterNode/MasterFileProcessing"
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/Transmission"
	"DistributedFileSystem/Utils"
	"errors"
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
	m.MasterServer.FileMetadataID[FileUUID] = &FileMetadata
	err = m.MasterServer.Logger.Log(FileMetadata.FileID, "Upload")
	if err != nil {
		return err
	}
	SendFileMetadataToMasterBackup(&FileMetadata, m.MasterServer)
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
	FileUUID := m.MasterServer.FileMetadataName[FileName].FileID

	err = m.MasterServer.Logger.Log(FileUUID, "Delete")
	if err != nil {
		return err
	}
	SendDeleteToMasterBackup(m.MasterServer.FileMetadataName[FileName], m.MasterServer)
	delete(m.MasterServer.FileMetadataName, FileName)
	delete(m.MasterServer.FileMetadataID, FileUUID)
	*reply = "Success!"
	return nil
}

func (m *MasterRPCServer) ReceiveSearchFromFrontendService(Filename string, reply *Metadata.FileMetaData) error {
	FileMetadata, exist := m.MasterServer.FileMetadataName[Filename]
	if !exist {
		*reply = Metadata.FileMetaData{}
		return nil
	}
	err := m.MasterServer.Logger.Log(FileMetadata.FileID, "Search")
	if err != nil {
		return err
	}
	SendSearchToMasterBackup(FileMetadata.FileID, m.MasterServer)
	*reply = *FileMetadata
	return nil
}

// Master Sync Log
func (m *MasterRPCServer) ReceiveSyncLogFromMasterBackup(LatestID int64, reply *bool) error {
	for i := int64(0); i < m.MasterServer.Logger.LatestID-LatestID; i++ {
		Log, err := m.MasterServer.Logger.ReadLog(i + LatestID + 1)
		if err != nil {
			return err
		}
		Type := Log[len(Log)-2:]
		if Type == "UP" {
			FileUUID := Log[30:65]
			FileMetadata := m.MasterServer.FileMetadataID[FileUUID]
			SendFileMetadataToMasterBackup(FileMetadata, m.MasterServer)
		} else if Type == "DE" {
			FileUUID := Log[30:65]
			FileMetadata := m.MasterServer.FileMetadataID[FileUUID]
			SendDeleteToMasterBackup(FileMetadata, m.MasterServer)
		} else {
			FileUUID := Log[30:65]
			SendSearchToMasterBackup(FileUUID, m.MasterServer)
		}
	}
	return nil
}

// MasterBackup receive Filemetada
func (m *MasterRPCServer) ReceiveFileMetadataFromMaster(FileMetadata *Metadata.FileMetaData, reply *bool) error {
	m.MasterServer.FileMetadataName[FileMetadata.FileName] = FileMetadata
	m.MasterServer.FileMetadataID[FileMetadata.FileID] = FileMetadata
	err := m.MasterServer.Logger.Log(FileMetadata.FileID, "Upload")
	if err != nil {
		return err
	}
	*reply = true
	return nil
}

func (m *MasterRPCServer) ReceiveSearchFromMaster(FileUUID string, reply *bool) error {
	err := m.MasterServer.Logger.Log(FileUUID, "Search")
	if err != nil {
		return err
	}
	*reply = true
	return nil
}

func (m *MasterRPCServer) ReceiveDeleteFromMaster(FileMetadata *Metadata.FileMetaData, reply *bool) error {
	_, exist := m.MasterServer.FileMetadataName[FileMetadata.FileName]
	if !exist {
		return errors.New("File does not exist in name mapping")
	}
	_, exist = m.MasterServer.FileMetadataID[FileMetadata.FileID]
	if !exist {
		return errors.New("File does not exist in ID mapping")
	}
	delete(m.MasterServer.FileMetadataName, FileMetadata.FileName)
	delete(m.MasterServer.FileMetadataID, FileMetadata.FileID)
	*reply = true
	return nil
}
