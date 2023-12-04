package main

import (
	"DistributedFileSystem/FrontendService"
	"DistributedFileSystem/LoadBalancer"
	"DistributedFileSystem/MasterNode"
	"DistributedFileSystem/SlaveNode"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(4)

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

	go func() {
		defer wg.Done()
		SlaveNode.StartSlaveServer()
	}()
	wg.Wait()
}
