package MasterFileProcessing

import (
	"DistributedFileSystem/MasterNode/MasterDefinition"
	"DistributedFileSystem/Transmission"
	"DistributedFileSystem/Utils"
	"os"
	"strconv"
)

func ChunkFile(FileArg *Transmission.FileArgs, m *MasterDefinition.MasterServer) ([]Transmission.ChunkArgs, error) {
	var FileChunks []Transmission.ChunkArgs
	err := os.WriteFile(m.StoragePath+FileArg.FileName, FileArg.Data, 0666)
	if err != nil {
		return make([]Transmission.ChunkArgs, 0), err
	}
	Shards, err := Utils.ShardFile(m.StoragePath+FileArg.FileName, m.ChunkSize)
	if err != nil {
		return make([]Transmission.ChunkArgs, 0), err
	}
	UUID, _ := Utils.GenerateUniqueID()
	for _, Shard := range Shards {
		ChunkID := UUID + "_" + strconv.Itoa(Shard.Index)
		Size := int64(len(Shard.Data))
		Data := Shard.Data
		temp := Transmission.ChunkArgs{ChunkName: ChunkID, Size: Size, Data: Data}
		FileChunks = append(FileChunks, temp)
	}
	return FileChunks, nil
}

func ReplicationChunks(ChunkArgs []Transmission.ChunkArgs, m *MasterDefinition.MasterServer) ([]Transmission.ChunkArgs, error) {
	var FileChunks []Transmission.ChunkArgs
	for i := 0; i < m.ReplicationFactor+1; i++ {
		for _, ChunkArg := range FileChunks {
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
