package main

import (
	"fmt"
	"sync"
)

const numberZero = 6      // Počet núl v hashi
var found uint32  // Indikátor, či bol nájdený výsledok
var offset uint64 // Offset, ktorý sa bude zväčšovať

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		startCPUWorkers(8)
	}()

	go func() {
		defer wg.Done()
		openCLWorker()
	}()

	wg.Wait()
	fmt.Println("Hľadanie ukončené.")
}
