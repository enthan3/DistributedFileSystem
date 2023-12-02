package main

import (
	"DistributedFileSystem/FrontendService"
	"DistributedFileSystem/LoadBalancer"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		FrontendService.StartFrontendServiceServer()

	}()
	go func() {
		defer wg.Done()
		LoadBalancer.StartLoadBalancerServer()
	}()

	wg.Wait()
	fmt.Printf("1")
}
