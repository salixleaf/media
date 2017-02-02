package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	mimage "github.com/salixleaf/media/image"
)

var debug bool
var clean bool

func main() {

	var image_path string
	flag.BoolVar(&debug, "debug", false, "specify debug mode")
	flag.BoolVar(&clean, "clean", false, "delete invalid image")
	flag.StringVar(&image_path, "dir", ".", "specify debug mode")
	help := flag.Bool("h", false, "print help page")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	runtime.GOMAXPROCS(2)
	ReadImageDir(image_path)
}

func ReadImageDir(picPath string) error {

	multi := NewMulti(4)
	err := filepath.Walk(picPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		multi.Incr()
		defer multi.Decr()

		p, _ := mimage.Stat(path)

		if p.Type == "" {
			if clean {
				fmt.Printf("delete invalid: %s\n", p.Path)
				os.Remove(p.Path)
			} else {
				fmt.Printf("invalid: %s\n", p.Path)
			}
		}
		if debug {
			fmt.Printf("Path: %s, width: %d, height: %d, size: %d, type: %s\n", p.Path, p.Width, p.Height, p.Size, p.Type)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() err: %v\n", err)
		return err
	}

	multi.Wait()
	return nil
}

type Multi struct {
	C chan bool
	W *sync.WaitGroup
}

func NewMulti(multi int) *Multi {
	return &Multi{
		C: make(chan bool, multi),
		W: new(sync.WaitGroup),
	}
}

func (this *Multi) Incr() {
	this.W.Add(1)
	this.C <- true
}

func (this *Multi) Decr() {
	this.W.Done()
	<-this.C
}

func (this *Multi) Wait() {
	this.W.Wait()
}
