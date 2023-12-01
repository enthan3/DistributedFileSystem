package FrontendServiceDefinition

import (
	"DistributedFileSystem/FrontendService/FrontendServiceCache"
)

// FrontendServiceServer Structure to store Master Address, Cache
type FrontendServiceServer struct {
	Master      string
	Cache       FrontendServiceCache.Cache
	StoragePath string
}
