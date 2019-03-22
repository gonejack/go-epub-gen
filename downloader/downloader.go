package downloader

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"log"
	"net/http"
	"os"
	"time"
)

var client = &grab.Client{
	HTTPClient: &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	},
	UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36",
	BufferSize: int(ResponseBuffer),
}
var concurrent = 5

func Batch(workers int, dstDir string, urls ...string) (<-chan *grab.Response, error) {
	fi, err := os.Stat(dstDir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("destination is not a directory")
	}

	reqs := make([]*grab.Request, len(urls))
	for i := 0; i < len(urls); i++ {
		req, err := grab.NewRequest(dstDir, urls[i])
		if err != nil {
			return nil, err
		}
		reqs[i] = req
		req.RateLimiter = limiter
	}

	ch := client.DoBatch(workers, reqs...)

	return ch, nil
}

func Download(dstDir string, urls ...string) (result map[string]string) {
	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	respch, err := Batch(concurrent, dstDir, urls...)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	completed := 0
	inProgress := 0
	responses := make([]*grab.Response, 0)
	result = make(map[string]string)
	ticker := time.NewTicker(time.Second)

	for completed < len(urls) {
		select {
		case resp, ok := <-respch:
			if ok && resp != nil {
				responses = append(responses, resp)
			}
		case <-ticker.C:
			// clear lines
			if inProgress > 0 {
				fmt.Printf("\033[%dA\033[K", inProgress)
			}

			// update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					if resp.Err() == nil {
						fmt.Printf("Finished %s %d / %d bytes (%d%%)\n", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
						result[resp.Request.URL().String()] = resp.Filename
					} else {
						_, _ = fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", resp.Request.URL(), resp.Err())
					}

					// mark completed
					responses[i] = nil
					completed++
				}
			}

			// update downloads in progress
			inProgress = 0
			for _, resp := range responses {
				if resp != nil {
					inProgress++
					fmt.Printf("Downloading %s %d / %d bytes (%d%%)\033[K\n", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
				}
			}
		}
	}

	ticker.Stop()

	fmt.Printf("%d files successfully downloaded.\n", len(urls))

	return
}