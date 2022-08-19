package gonmap

import (
	"context"
	"goon3/lib/gonmap/simplenet"
	"net"
	"strings"
	"time"
)

func PortScan(protocol string, netloc string, duration time.Duration) bool {
	result, err := simplenet.Send(protocol, netloc, "", duration, 0)
	if err == nil {
		return true
	}
	if len(result) > 0 {
		return true
	}
	if strings.Contains(err.Error(), "STEP1") {
		return false
	} else {
		return true
	}
}

func DnsScan(addr string) bool {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 7 * time.Second,
			}
			return d.DialContext(ctx, "udp", addr)
		},
	}

	_, err := r.LookupHost(context.Background(), "localhost")
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return false
		}
	}
	return true
}
