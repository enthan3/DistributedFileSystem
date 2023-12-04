package MasterFileProcessing

import (
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/Transmission"
	"DistributedFileSystem/Utils"
	"os"
	"strconv"
)

func ChunkFile(FileUUID string, FileArg *Transmission.FileArgs, m *MasterDefinition.MasterServer) ([]Transmission.ChunkArg, error) {
	var FileChunks []Transmission.ChunkArg
	err := os.WriteFile(m.StoragePath+FileArg.FileName, FileArg.Data, 0666)
	if err != nil {
		return make([]Transmission.ChunkArg, 0), err
	}
	Shards, err := Utils.ShardFile(m.StoragePath+FileArg.FileName, m.ChunkSize)
	if err != nil {
		return make([]Transmission.ChunkArg, 0), err
	}
	for _, Shard := range Shards {
		ChunkID := FileUUID + "_" + strconv.Itoa(Shard.Index)
		Size := int64(len(Shard.Data))
		Data := Shard.Data
		temp := Transmission.ChunkArg{ChunkName: ChunkID, Size: Size, Data: Data}
		FileChunks = append(FileChunks, temp)
	}
	err = os.Remove(m.StoragePath + FileArg.FileName)
	if err != nil {
		return make([]Transmission.ChunkArg, 0), err
	}
	return FileChunks, nil
}

func ReplicationChunks(ChunkArgs []Transmission.ChunkArg, m *MasterDefinition.MasterServer) ([]Transmission.ChunkArg, error) {
	var FileChunks []Transmission.ChunkArg
	for i := 0; i < m.ReplicationFactor+1; i++ {
		for _, ChunkArg := range ChunkArgs {
			if i == 0 {
				ChunkArg.ChunkName += "_" + "main"
				FileChunks = append(FileChunks, ChunkArg)
			} else {
				ChunkArg.ChunkName += "_" + "replicate" + strconv.Itoa(i)
				FileChunks = append(FileChunks, ChunkArg)
			}
		}
	}
	return FileChunks, nil
}
