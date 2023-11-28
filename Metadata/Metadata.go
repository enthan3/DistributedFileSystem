package Metadata

import "time"

type FileMetaData struct {
	FileID         string
	FileName       string
	Size           int64
	CreateTime     time.Time
	ChunkMain      []ChunkMetaData
	ChunkReplicate []ChunkMetaData
}

type ChunkMetaData struct {
	ChunkName  string
	Size       int64
	CreateTime time.Time
	ChunkNode  string
}
