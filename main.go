package main

import (
	"DistributedFileSystem/FrontendService"
	"DistributedFileSystem/LoadBalancer"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2) // 为两个 goroutine 增加计数器

	go func() {
		defer wg.Done() // 在 goroutine 完成时调用 Done()
		FrontendService.StartFrontendServiceServer()

	}()
	go func() {
		defer wg.Done()
		LoadBalancer.StartLoadBalancerServer()
	}()

	wg.Wait() // 阻塞，直到所有 goroutine 调用 Done()
}
