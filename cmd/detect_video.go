package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/salixleaf/media/video"
)

func main() {

	debug := flag.Bool("debug", false, "specify debug mode")
	path := flag.String("dir", ".", "specify video's directory")
	tidy := flag.Bool("tidy", false, "tidy video with resolution")
	help := flag.Bool("h", false, "print help page")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	filepath.Walk(*path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		v := video.Stat(path)

		var height int64
		if v.Height > v.Width {
			height = v.Width
		} else {
			height = v.Height
		}

		var label string
		switch {
		case height >= 1080:
			label = "1080"
		case height >= 720:
			label = "720"
		case height >= 640:
			label = "640"
		case height >= 480:
			label = "480"
		case height >= 320:
			label = "320"
		default:
			label = "oth"
		}

		if *tidy {
			MoveFile(path, dir+"."+label)
		}

		if *debug {
			fmt.Printf("Path: %s, width: %d, height: %d, bitrate: %.2fkb/s, length: %.2fs, size: %dByte\n", v.Path, v.Width, v.Height, v.Bitrate, v.Length, v.Size)
		}
		return nil
	})
}

func MoveFile(file string, dstPath string) {
	fmt.Println(file, dstPath)
	if !FileExist(dstPath) {
		os.Mkdir(dstPath, 0755)
	}

	os.Rename(file, filepath.Join(dstPath, filepath.Base(file)))
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
