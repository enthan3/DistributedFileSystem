package FrontendServiceRPC

import (
	"DistributedFileSystem/FrontendService/FrontendServiceDefinition"
	"DistributedFileSystem/Metadata"
	"DistributedFileSystem/Transmission"
	"errors"
	"net/rpc"
)

func SendFileToMaster(FileArg *Transmission.FileArgs, f *FrontendServiceDefinition.FrontendServiceServer) error {
	var reply bool
	client, err := rpc.Dial("tcp", f.Master)
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

func SendDeleteToMaster(FileName string, f *FrontendServiceDefinition.FrontendServiceServer) error {
	var reply string
	client, err := rpc.Dial("tcp", f.Master)
	if err != nil {
		return err
	}
	//TODO
	err = client.Call("", FileName, &reply)
	if err != nil {
		return err
	}
	if reply == "File does not exist" {
		return errors.New(reply)
	}
	return nil
}

func SendSearchToMaster(FileName string, f *FrontendServiceDefinition.FrontendServiceServer) (Metadata.FileMetaData, error) {
	var reply Metadata.FileMetaData
	client, err := rpc.Dial("tcp", f.Master)
	if err != nil {
		return Metadata.FileMetaData{}, err
	}
	err = client.Call("", FileName, &reply)
	if err != nil {
		return Metadata.FileMetaData{}, err
	}
	if reply.FileName == "" && reply.FileID == "" {
		return Metadata.FileMetaData{}, errors.New("File does not exist")
	}
	return reply, nil
}

func SendDownloadToSlave(chunkName string, Slave string, f *FrontendServiceDefinition.FrontendServiceServer) (Transmission.ChunkArgs, error) {
	var reply Transmission.ChunkArgs
	client, err := rpc.Dial("tcp", Slave)
	if err != nil {
		return Transmission.ChunkArgs{}, err
	}
	err = client.Call("", chunkName, &reply)
	if err != nil {
		return Transmission.ChunkArgs{}, err
	}
	return reply, nil
}
