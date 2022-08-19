package gonmap

import (
	"context"
	"goon3/lib/kscan/lib/urlparse"
	"time"
)

func GetTcpBanner(netloc string, nmap *Nmap, timeout time.Duration) *TcpBanner {
	parse, _ := urlparse.Load(netloc)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resChan := make(chan *TcpBanner)
	go func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		r := nmap.Scan(parse.Netloc, parse.Port)
		resChan <- &r
	}()

	for {
		select {
		case <-ctx.Done():
			close(resChan)
			return nil
		case res := <-resChan:
			close(resChan)
			return res
		}
	}
}
