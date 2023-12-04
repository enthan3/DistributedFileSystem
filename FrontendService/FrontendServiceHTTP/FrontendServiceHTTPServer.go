package FrontendServiceHTTP

import (
	"DistributedFileSystem/FrontendService/FrontendServiceDefinition"
	"DistributedFileSystem/FrontendService/FrontendServiceRPC"
	"DistributedFileSystem/Transmission"
	"DistributedFileSystem/Utils"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

//加入文件存在返回文件metaData加入不存在返回false

// ReceiveFileFromFrontend Receive and store the file temporarily from frontend
func ReceiveFileFromFrontend(w http.ResponseWriter, r *http.Request, f *FrontendServiceDefinition.FrontendServiceServer) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	file, handler, err := r.FormFile("File")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	FileArg := Transmission.FileArgs{
		FileName: handler.Filename,
		Size:     handler.Size,
		Data:     data,
	}
	err = FrontendServiceRPC.SendFileToMaster(&FileArg, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ReceiveDeleteFromFrontend Receive delete request and delete metadata from cache and master node, real data from slave nodes
func ReceiveDeleteFromFrontend(w http.ResponseWriter, r *http.Request, f *FrontendServiceDefinition.FrontendServiceServer) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	fileName := r.URL.Query().Get("File")
	err := FrontendServiceRPC.SendDeleteToMaster(fileName, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f.Cache.Del(fileName)
	return
}

// ReceiveDownloadFromFrontend Receive downloading request and retrieve related chunks from slave nodes and form file from chunks
func ReceiveDownloadFromFrontend(w http.ResponseWriter, r *http.Request, f *FrontendServiceDefinition.FrontendServiceServer) {
	var Shards []Utils.Shard
	if r.Method != "GET" {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	FileName := r.URL.Query().Get("File")
	FileMetaData, exist := f.Cache.Get(FileName)
	if !exist {
		tmp, err := FrontendServiceRPC.SendSearchToMaster(FileName, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		FileMetaData = &tmp
		f.Cache.Put(FileMetaData)
	}
	Shards, err := FrontendServiceRPC.SendDownloadToSlaves(FileMetaData, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = Utils.FormFile(Shards, FileMetaData.FileName, f.StoragePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment;filename="+FileMetaData.FileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	FileData, err := os.ReadFile(f.StoragePath + FileMetaData.FileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(FileData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = os.Remove(f.StoragePath + FileMetaData.FileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ReceiveSearchFromFrontend Receive search request and return file metadata if exist
func ReceiveSearchFromFrontend(w http.ResponseWriter, r *http.Request, f *FrontendServiceDefinition.FrontendServiceServer) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	FileName := r.URL.Query().Get("File")

	FileMetaData, exist := f.Cache.Get(FileName)
	if !exist {
		temp, err := FrontendServiceRPC.SendSearchToMaster(FileName, f)
		if err != nil {
			return
		}
		f.Cache.Put(&temp)
		FileMetaData = &temp
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(FileMetaData)
	if err != nil {
		http.Error(w, "Error encoding search results", http.StatusInternalServerError)
		return
	}
}
