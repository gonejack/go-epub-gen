package main

import (
	"fmt"
	"go-epub-gen/downloader"
	"os"
	"path/filepath"
)

func main() {
	// get URL to download from command args
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s url [url]...\n", os.Args[0])
		os.Exit(1)
	}

	urls := os.Args[1:]


	downloader.Download(filepath.Join(os.TempDir(), "go-epub-gen"), urls...)
}