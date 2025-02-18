package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	sha256simd "github.com/minio/sha256-simd" // NEON optimalizovaná SHA-256
)

const prefix = "para"
const batchSize = 1024*1024*8 // Každý worker spracuje 100 čísel

func cpuWorker(id int) {
	for atomic.LoadUint32(&found) == 0 {
		localStart := uint64(atomic.AddUint64(&offset, batchSize)) - uint64(batchSize) // Pôvodný offset pre tento batch
		fmt.Println("cpu ", id, " add")
		for i := uint64(0); i < batchSize; i++ {
			num := localStart + i
			word := fmt.Sprintf("%s%d", prefix, num)

			hash := sha256Compute(word)

			if hasLeadingZeros(hash, numberZero) {
				atomic.StoreUint32(&found, 1)
				fmt.Printf("[CPU-%d] Found! Number: %d, Hash: %s\n", id, num, hex.EncodeToString(hash[:]))
				return
			}
		}
	}
}

func sha256Compute(word string) [32]byte {
	if runtime.GOARCH == "arm64" {
		return sha256simd.Sum256([]byte(word)) // NEON optimalizovaná SHA-256
	}
	return sha256.Sum256([]byte(word)) // Normálna Go SHA-256
}

func hasLeadingZeros(hash [32]byte, numZero int) bool {
	zeroCount := 0

	for _, b := range hash {
		if b == 0 {
			zeroCount++ // Celý bajt je nulový
			if zeroCount >= numZero {
				return true
			}
		} else {
			return false // Ak narazíme na nenulový bajt, zastavíme kontrolu
		}
	}

	return false
}

func startCPUWorkers(numWorkers int) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(ii int) {
			fmt.Println("start cpu", ii)
			cpuWorker(ii)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
