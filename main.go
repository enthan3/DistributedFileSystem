package main

import (
	"DistributedFileSystem/FrontendService"
	"DistributedFileSystem/LoadBalancer"
	"DistributedFileSystem/MasterNode"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		FrontendService.StartFrontendServiceServer()

	}()
	go func() {
		defer wg.Done()
		LoadBalancer.StartLoadBalancerServer()
	}()

	go func() {
		defer wg.Done()
		MasterNode.StartMasterServer()
	}()
	wg.Wait()
	fmt.Printf("1")
}
