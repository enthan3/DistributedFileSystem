package main

import (
	"DistributedFileSystem/MasterNode"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		MasterNode.StartMasterServer()
	}()

	wg.Wait()
}
