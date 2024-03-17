package main

import (
	"fmt"
	"github.com/opencl-pure/constantsCL"
	opencl "github.com/opencl-pure/highCL"
	pure "github.com/opencl-pure/pureCL"
	"image"
	"log"
	"strings"
)

func main() {
	err := opencl.Init(pure.Version2_0) //init with version of OpenCL
	if err != nil {
		log.Fatal(err)
	}
	d, err := opencl.GetDefaultDevice()
	if err != nil {
		log.Fatal(err)
	}
	defer func(d *opencl.Device) {
		err = d.Release()
		if err != nil {
			log.Println(err)
		}
	}(d)
	invertedImg, err := d.NewImage2D(constantsCL.CL_RGBA, image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func(fractalImg *opencl.Image) {
		err = fractalImg.Release()
		if err != nil {
			log.Println(err)
		}
	}(invertedImg)
	var bigKernelBuilder strings.Builder
	for i := 0; i < len(kernels); i++ {
		data, _ := getKernelData(kernels[i])
		bigKernelBuilder.WriteString(fmt.Sprintln(data))
		//time.Sleep(time.Second)
	}
	_, err = d.AddProgram(bigKernelBuilder.String())
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(kernels); i++ {
		makeOneImg(d, invertedImg, kernels[i])
		//time.Sleep(time.Second)
	}
}
