package Utils

import (
	"io"
	"os"
	"sort"
)

type Shard struct {
	Index int
	Data  []byte
}

func FormFile(shards []Shard, fileName string, StoragePath string) error {
	file, err := os.Create(StoragePath + fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	sort.Slice(shards, func(i, j int) bool {
		return shards[i].Index < shards[j].Index
	})

	for i := 0; i < len(shards); i++ {
		_, err := file.Write(shards[i].Data)
		if err != nil {
			return err
		}
	}
	return err
}

func ShardFile(fileName string, shardSize int64) ([]Shard, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var shards []Shard
	fileSize := fileInfo.Size()

	for i := int64(0); i < fileSize; i += shardSize {
		buffer := make([]byte, shardSize)
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		shard := Shard{
			Index: int(i / shardSize),
			Data:  buffer[:bytesRead],
		}
		shards = append(shards, shard)

		if err == io.EOF {
			break
		}
	}
	return shards, nil
}
