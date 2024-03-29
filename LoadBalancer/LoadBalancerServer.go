package LoadBalancer

import (
	"DistributedFileSystem/LoadBalancer/LoadBalancerConfiguration"
	"DistributedFileSystem/LoadBalancer/LoadBalancerDefinition"
	"DistributedFileSystem/LoadBalancer/LoadBalancerHTTP"
	"DistributedFileSystem/LoadBalancer/LoadBalancerRPC"
	"DistributedFileSystem/Transmission"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func StartLoadBalancerServer() {
	Config, err := LoadBalancerConfiguration.LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	l := LoadBalancerRPC.LoadBalancerRPCServer{LoadBalancerServer: &LoadBalancerDefinition.LoadBalancerServer{MasterStatusMap: make(map[string]*Transmission.MasterStatusArg),
		MasterBackupsMap: Config.MastersAddress, CurrentHTTP: Config.CurrentHTTPAddress, CurrentRPC: Config.CurrentRPCAddress, ServiceHTTP: Config.ServiceHTTPAddress, ServiceRPC: Config.ServiceRPCAddress}}
	for MasterAddress, _ := range l.LoadBalancerServer.MasterBackupsMap {
		l.LoadBalancerServer.MasterStatusMap[MasterAddress] = &Transmission.MasterStatusArg{}
	}
	err = rpc.Register(&l)
	if err != nil {
		log.Fatal(err)
	}
	L, err := net.Listen("tcp", Config.CurrentRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	err = LoadBalancerRPC.SendMastersToFrontendService(l.LoadBalancerServer)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err = LoadBalancerRPC.SendStatusRequestToMaster(l.LoadBalancerServer)
			fmt.Println(l.LoadBalancerServer.MasterStatusMap["192.168.8.15:10010"].Stat.Load1)
			if err != nil {
				log.Fatal(err)
			}
			err = LoadBalancerRPC.SendLowestUsageMasterToFrontendService(l.LoadBalancerServer)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(Config.HeartbeatDuration) * time.Second)
		}
	}()

	go func() {
		for {
			conn, err := L.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go rpc.ServeConn(conn)

		}
	}()
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			LoadBalancerHTTP.RedirectToFrontendService(w, r, l.LoadBalancerServer)
		})

		err = http.ListenAndServe(l.LoadBalancerServer.CurrentHTTP, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
