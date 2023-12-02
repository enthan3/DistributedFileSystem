package FrontendServiceDefinition

import (
	"DistributedFileSystem/FrontendService/FrontendServiceCache"
)

// FrontendServiceServer Structure to store Master Address, Cache
type FrontendServiceServer struct {
	//Map structure Master Address to Master Backup
	Masters            []string
	LowestMaster       string
	Cache              *FrontendServiceCache.Cache
	StoragePath        string
	CurrentHTTPAddress string
	CurrentRPCAddress  string
}
