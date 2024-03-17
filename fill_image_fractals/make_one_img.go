package main

import (
	"fmt"
	opencl "github.com/opencl-pure/highCL"
	"image/png"
	"log"
	"os"
)

func makeOneImg(d *opencl.Device, fractalImg *opencl.Image, kernelName string) {
	k, err := d.Kernel(kernelName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(k *opencl.Kernel) {
		err = k.ReleaseKernel()
		if err != nil {
			log.Println(err)
		}
	}(k)
	event, err := k.Global(width, height).Local(1, 1).Run(nil, fractalImg)
	if err != nil {
		log.Fatal(err)
	}
	defer func(event *opencl.Event) {
		err = event.Release()
		if err != nil {
			log.Println(err)
		}
	}(event)
	_ = event.Wait()
	_ = k.Flush()
	_ = k.Finish()
	fractal, err := fractalImg.Data()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("outputs/" + kernelName + "_fractal.png")
	if err != nil {
		log.Fatal(err)
	}
	_ = png.Encode(f, fractal)
	_ = f.Close()
}

func getKernelData(name string) (string, error) {
	b, err := os.ReadFile("kernels/" + name + ".cl")
	if err != nil {
		return "", err
	}
	return fmt.Sprint(string(b)), nil
}
