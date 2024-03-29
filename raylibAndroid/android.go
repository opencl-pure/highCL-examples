//go:build android

package main

import (
	"embed"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"strings"
	"syscall"
)

func doSpecific() {
	cErr := make([]string, 0, 2)
	env := []string{"LD_LIBRARY_PATH=/system/vendor/lib64"}
	dir := strings.Replace(rl.HomeDir(), "/user/0/", "/data/", 1)
	file := dir + "/wrap.sh"
	if ds, err := os.ReadDir(dir); err == nil {
		if !find(ds, "wrap.sh") {
			err = move("wrap.sh", dir)
			if err != nil {
				rl.TraceLog(3, err.Error())
			}
			err = move("fill_image_fractals", dir)
			if err != nil {
				rl.TraceLog(3, err.Error())
			}
			fillKernels(dir)
		}
	} else {
		rl.TraceLog(3, err.Error())
	}

	forkExec, err := syscall.ForkExec("/system/bin/chmod",
		[]string{"chmod", "1777", file},
		&syscall.ProcAttr{
			Dir:   dir,
			Env:   env,
			Files: nil,
			Sys:   nil,
		})
	if err != nil {
		cErr = append(cErr, err.Error())
		goto print
	}
	_, err = syscall.Wait4(forkExec, nil, 0, nil)
	if err != nil {
		cErr = append(cErr, err.Error())
	}
	forkExec, err = syscall.ForkExec("/system/bin/su",
		[]string{"su", "-c", file},
		&syscall.ProcAttr{
			Dir:   dir,
			Env:   env,
			Files: nil,
			Sys:   nil,
		})
	if err != nil {
		cErr = append(cErr, err.Error())
		goto print
	}
	_, err = syscall.Wait4(forkExec, nil, 0, nil)
	if err != nil {
		cErr = append(cErr, err.Error())
	}
print:
	rl.TraceLog(3, fmt.Sprint(cErr))
	if len(cErr) == 0 {
		loadKernelImages()
	}
}

//go:embed kernels
var ker embed.FS

func fillKernels(dir string) {
	dirKernels := dir + "/kernels"
	forkExec, err := syscall.ForkExec("/system/bin/mkdir",
		[]string{"mkdir", dirKernels},
		&syscall.ProcAttr{
			Dir:   dir,
			Env:   nil,
			Files: nil,
			Sys:   nil,
		})
	if err != nil {
		rl.TraceLog(3, err.Error())
	}
	_, err = syscall.Wait4(forkExec, nil, 0, nil)
	for _, kernel := range kernels {
		k := kernel + ".cl"
		file, err := ker.ReadFile("kernels/" + k)
		if err != nil {
			rl.TraceLog(3, err.Error())
			continue
		}
		err = os.WriteFile(dirKernels+"/"+k, file, 1777)
		if err != nil {
			rl.TraceLog(3, err.Error())
		}
	}
}

func find(ds []os.DirEntry, s string) bool {
	if ds == nil || len(ds) == 0 {
		return false
	}
	for _, dirEntry := range ds {
		if dirEntry.Name() == s {
			return true
		}
	}
	return false
}

func move(name, dir string) error {
	asset, err := rl.OpenAsset(name)
	if err != nil {
		return err
	}
	b, bb := make([]byte, 1024), make([]byte, 0, 1024*1024*3)
	for {
		n0, err0 := asset.Read(b)
		if err0 != nil {
			err = err0
			break
		}
		bb = append(bb, b[:n0]...)
	}
	if len(bb) == 0 {
		return err
	}
	return os.WriteFile(dir+"/"+name, bb, 1777)
}
