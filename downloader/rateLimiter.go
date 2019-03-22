package downloader

import (
	"context"
	"github.com/inhies/go-bytesize"
	"time"
)

var BytesPerSecond = 1 * bytesize.MB
var ResponseBuffer = 16 * bytesize.KB

var limiter *RateLimiter

func init() {
	limiter = &RateLimiter{
		ch: make(chan bytesize.ByteSize),
	}

	limiter.start()
}

type RateLimiter struct {
	ch chan bytesize.ByteSize
}

func (lim *RateLimiter) start() {
	var ticker = make(chan bytesize.ByteSize)

	go func() {
		var available = BytesPerSecond / 100
		for {
			ticker <- available

			time.Sleep(time.Second / 100)
		}
	}()

	go func() {
		var available bytesize.ByteSize
		for {
			for available += <-ticker; available > ResponseBuffer; {
				lim.ch <- ResponseBuffer

				available -= ResponseBuffer
			}
		}
	}()
}

func (lim *RateLimiter) WaitN(ctx context.Context, n int) (err error) {
	var got int

	for got < n {
		if err = ctx.Err(); err == nil {
			got += int(<- lim.ch)
		} else {
			return
		}
	}

	return
}