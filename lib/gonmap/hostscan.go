package gonmap

import (
	"fmt"
	"github.com/go-ping/ping"
	"goon3/lib/kscan/lib/slog"
	"runtime"
	"time"
)

func HostDiscovery(ip string) (online bool) {
	online = pingCheck(ip)
	if online {
		return true
	}
	online = tcpCheck(ip)
	if online {
		return true
	}
	return false
}

func HostDiscoveryIcmp(ip string) (online bool) {
	online = pingCheck(ip)
	if online {
		return true
	}
	return false
}

func pingCheck(ip string) bool {
	p, err := ping.NewPinger(ip)
	if runtime.GOOS == "windows" {
		p.SetPrivileged(true)
	}
	if err != nil {
		slog.Debug(err.Error())
		return false
	}
	p.Count = 2
	p.Timeout = time.Second * 2
	err = p.Run() // Blocks until finished.
	if err != nil {
		slog.Debug(err.Error())
	}
	s := p.Statistics()
	if s.PacketsRecv > 0 {
		return true
	}
	return false
}

func tcpCheck(ip string) bool {
	tcpArr := []int{21, 22, 23, 80, 443, 445, 8080, 3389}
	for _, port := range tcpArr {
		netloc := fmt.Sprintf("%s:%d", ip, port)
		online := PortScan("tcp", netloc, time.Second*2)
		if online {
			return true
		}
	}
	return false
}
