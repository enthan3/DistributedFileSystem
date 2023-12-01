package FrontendServiceHTTP

import (
	"DistributedFileSystem/FrontendService/FrontendServiceDefinition"
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
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Retrieving file FrontendService error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Read file FrontendService error", http.StatusInternalServerError)
		return
	}
	err = os.WriteFile(handler.Filename, data, 0666)
	if err != nil {
		http.Error(w, "Downloading file FrontendService error", http.StatusInternalServerError)
		return
	}

}

// ReceiveDeleteFromFrontend Receive delete request and delete metadata from cache and master node, real data from slave nodes
func ReceiveDeleteFromFrontend(w http.ResponseWriter, r *http.Request, f *FrontendServiceDefinition.FrontendServiceServer) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	fileName := r.URL.Query().Get("filename")
	//TODO （主节点）实现删除文件功能

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
	FileName := r.URL.Query().Get("filename")
	FileMetaData, exist := f.Cache.Get(FileName)
	if !exist {
		//TODO （主节点）实现下载文件返回文件元数据
		f.Cache.Put(FileMetaData)
	}
	//TODO （从节点）实现下载文件返回文件真实数据

	err := Utils.FormFile(Shards, FileMetaData.FileName, f.StoragePath)
	if err != nil {
		http.Error(w, "Form file FrontendService error", http.StatusInternalServerError)
		return
	}
	defer os.Remove(f.StoragePath + FileMetaData.FileName)
	w.Header().Set("Content-Disposition", "attachment;filename="+FileMetaData.FileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	FileData, err := os.ReadFile(f.StoragePath + FileMetaData.FileName)
	if err != nil {
		http.Error(w, "Reading file FrontendService error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(FileData)
	if err != nil {
		http.Error(w, "Write write FrontendService error", http.StatusInternalServerError)
		return
	}
}

// ReceiveSearchFromFrontend Receive search request and return file metadata if exist
func ReceiveSearchFromFrontend(w http.ResponseWriter, r *http.Request, f *FrontendServiceDefinition.FrontendServiceServer) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	FileName := r.URL.Query().Get("filename")

	FileMetaData, exist := f.Cache.Get(FileName)
	if !exist {
		//TODO 实现SendSearchToMaster方法，文件存在返回文件元数据不存在返回false

		f.Cache.Put(FileMetaData)
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(FileMetaData)
	if err != nil {
		http.Error(w, "Error encoding search results", http.StatusInternalServerError)
		return
	}
}
