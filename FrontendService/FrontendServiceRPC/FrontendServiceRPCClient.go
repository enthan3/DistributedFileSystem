package FrontendServiceRPC

import (
	"DistributedFileSystem/FrontendService/FrontendServiceDefinition"
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/Transmission"
	"DistributedFileSystem/Utils"
	"errors"
	"net/rpc"
)

//If file does not exist it will not return error, it will reply false

// SendFileToMaster Send file to the lowest usage master
func SendFileToMaster(FileArg *Transmission.FileArgs, f *FrontendServiceDefinition.FrontendServiceServer) error {
	var reply bool
	client, err := rpc.Dial("tcp", f.LowestMaster)
	if err != nil {
		return err
	}
	//TODO
	err = client.Call("", FileArg, &reply)
	if err != nil {
		return err
	}
	if reply != true {
		return errors.New("Receive file at Master error")
	}
	return nil
}

// SendDeleteToMaster send delete request to every master to delete the file
func SendDeleteToMaster(FileName string, f *FrontendServiceDefinition.FrontendServiceServer) error {
	var reply string
	for _, MasterAddress := range f.Masters {
		client, err := rpc.Dial("tcp", MasterAddress)
		if err != nil {
			return err
		}
		//TODO
		err = client.Call("", FileName, &reply)
		if err != nil {
			return err
		}
		if reply == "Success!" {
			return nil
		}
	}
	return errors.New("File does not exist!")
}

func SendSearchToMaster(FileName string, f *FrontendServiceDefinition.FrontendServiceServer) (Metadata.FileMetaData, error) {
	var reply Metadata.FileMetaData
	for _, MasterAddress := range f.Masters {
		client, err := rpc.Dial("tcp", MasterAddress)
		if err != nil {
			return Metadata.FileMetaData{}, err
		}
		//TODO
		err = client.Call("", FileName, &reply)
		if err != nil {
			return Metadata.FileMetaData{}, err
		}
		if reply.FileName != "" && reply.FileID != "" {
			return reply, nil
		}
	}
	return Metadata.FileMetaData{}, errors.New("File does not exist")
}

func SendDownloadToSlaves(FileMetadata *Metadata.FileMetaData, f *FrontendServiceDefinition.FrontendServiceServer) ([]Utils.Shard, error) {
	var reply Transmission.ChunkArgs
	var Shards []Utils.Shard
	for index, ChunkMetadata := range FileMetadata.ChunkMain {
		client, err := rpc.Dial("tcp", ChunkMetadata.ChunkNode)
		if err != nil {
			return make([]Utils.Shard, 0), err
		}
		//TODO
		err = client.Call("", ChunkMetadata.ChunkName, &reply)
		if err != nil {
			return make([]Utils.Shard, 0), err
		}
		if reply.ChunkName == "" {
			return make([]Utils.Shard, 0), err
		}
		Shards = append(Shards, Utils.Shard{Index: index, Data: reply.Data})
	}
	return Shards, nil
}
