package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/opencl-pure/highCL"
	pure "github.com/opencl-pure/pureCL"
)
const numWorkItems = 4096*128 // Počet paralelných prác
const stepPerItem = 4096    // kolko Každý worker skúša  čísel

func openCLWorker() {
	runtime.LockOSThread() // OpenCL nesmie meniť thread
	defer runtime.UnlockOSThread()
	// OpenCL inicializácia
	err := highCL.Init(pure.Version2_0, "libOpenCL.so")
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}

	device, err := highCL.GetDefaultDevice()
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}
	defer device.Release()

	kernelSource, err := os.ReadFile("kernel.cl")
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}

	_, err = device.AddProgram(string(kernelSource))
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}

	// Nastavenie nemeniacich sa hodnôt
	prefix := []byte{'p', 'a', 'r', 'a', 0}
	prefixBuf, err := device.NewVector(prefix)
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}
	defer prefixBuf.Release()

	zerosBuf, err := device.NewVector([]uint8{numberZero})
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}
	defer zerosBuf.Release()
	foundBuf, err := device.NewVector([]int{-1})
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}
	defer foundBuf.Release()

	temp := make([]uint32, 10*4)
	resultBuf, err := device.NewVector(temp)
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}
	defer resultBuf.Release()

	offsetBuf, err := device.NewVector([]uint64{0, uint64(stepPerItem)})
	defer offsetBuf.Release()

	k, err := device.Kernel("hashKernel")
	if err != nil {
		log.Fatalf("Failed to add program: %v", err)
		return
	}
	defer k.ReleaseKernel()
	fmt.Println("starting gpu")
	for atomic.LoadUint32(&found) == 0 {
		currentOffset := atomic.AddUint64(&offset, uint64(numWorkItems*stepPerItem)) - uint64(numWorkItems*stepPerItem)
		fmt.Println("kernel add")

		err := <-offsetBuf.Reset([]uint64{currentOffset, uint64(stepPerItem)})
		if err != nil {
			log.Fatalf("Failed to write offset buffer: %v", err)
		}
		event, err := k.Global(numWorkItems).Local(1).Run(nil, prefixBuf, offsetBuf, resultBuf, zerosBuf, foundBuf)
		if err != nil {
			log.Fatalf("Failed to run kernel: %v", err)
			return
		}
		k.Flush()
		k.Finish()
		event.Release()

		result, err := foundBuf.Data()
		if err != nil {
			log.Fatalf("Failed to read result buffer: %v", err)
		}
		if int(result.Index(0).Int()) != -1 {
			h, err := resultBuf.Data()

			if err != nil {
				log.Fatalf("Failed to read result buffer: %v", err)
			}
			if atomic.CompareAndSwapUint32(&found, 0, 1) {
				candidate := (uint64(result.Index(0).Int()))
				fmt.Printf("[GPU] Found candidate: %d\nHash: ", candidate)
				for i := 2; i < 10; i++ {
					fmt.Printf("%08x", uint32(h.Index(i).Uint()))
				}
				fmt.Println()
			}
			return
		}
	}

}
