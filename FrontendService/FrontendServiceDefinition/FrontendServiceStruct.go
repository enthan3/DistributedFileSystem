package FrontendServiceDefinition

import "DistributedFileSystem/Metadata"

// FrontendServiceServer Structure to store Master Address, Cache and CacheCount
type FrontendServiceServer struct {
	Master     string
	Cache      map[string]Metadata.FileMetaData
	CacheCount int
}
