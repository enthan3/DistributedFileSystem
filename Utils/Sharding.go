package Utils

import (
	"os"
	"sort"
)

type Shard struct {
	Index int
	Data  []byte
}

// FormFile 将所有的分片拼装在一起,或者将单个分块的文件变成一个file
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
