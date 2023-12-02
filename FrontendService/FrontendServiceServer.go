package FrontendService

import (
	"DistributedFileSystem/FrontendService/FrontendServiceCache"
	"DistributedFileSystem/FrontendService/FrontendServiceConfiguration"
	"DistributedFileSystem/FrontendService/FrontendServiceDefinition"
	"DistributedFileSystem/FrontendService/FrontendServiceHTTP"
	"DistributedFileSystem/FrontendService/FrontendServiceRPC"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

func StartFrontendServiceServer() {
	Config, err := FrontendServiceConfiguration.LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var Cache FrontendServiceCache.Cache
	f := FrontendServiceRPC.FrontendServiceRPCServer{FrontendServiceServer: &FrontendServiceDefinition.FrontendServiceServer{Cache: Cache.NewCache(Config.CacheSize), StoragePath: Config.StoragePath, CurrentHTTPAddress: Config.HTTPAddress, CurrentRPCAddress: Config.RPCAddress}}
	err = rpc.Register(&f)
	if err != nil {
		log.Fatal(err)
	}
	L, err := net.Listen("tcp", f.FrontendServiceServer.CurrentRPCAddress)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			conn, err := L.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go rpc.ServeConn(conn)
			print(f.FrontendServiceServer.Masters)
		}
	}()

	go func() {
		defer wg.Done()
		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			FrontendServiceHTTP.ReceiveFileFromFrontend(w, r, f.FrontendServiceServer)
		})

		http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
			FrontendServiceHTTP.ReceiveDeleteFromFrontend(w, r, f.FrontendServiceServer)
		})

		http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
			FrontendServiceHTTP.ReceiveSearchFromFrontend(w, r, f.FrontendServiceServer)
		})

		http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
			FrontendServiceHTTP.ReceiveDownloadFromFrontend(w, r, f.FrontendServiceServer)
		})

		err = http.ListenAndServe(f.FrontendServiceServer.CurrentHTTPAddress, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

}
