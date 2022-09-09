package Host

import (
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"sync"
	"time"
)

const (
	REACHABLE   = "Yes"     //可达
	UNREACHABLE = "No"      //不可达
	UNKNOWN     = "Unknown" //未知
)

type Host struct {
	Ipaddr    string
	Reachable string
}

func NewHost(ipaddr string) *Host {
	return &Host{
		Ipaddr:    ipaddr,
		Reachable: UNKNOWN,
	}
}

func (h *Host) Ping(n int, tcp bool) {
	pinger, err := ping.NewPinger(h.Ipaddr)
	if err != nil {
		log.Println(err)
	}
	pinger.Count = n
	pinger.Timeout = time.Second
	pinger.SetPrivileged(tcp)

	// 运行pinger
	err = pinger.Run()
	if err != nil {
		fmt.Println(err)
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		h.Reachable = REACHABLE
	} else {
		h.Reachable = UNREACHABLE
	}
}

func (h *Host) Run(wg *sync.WaitGroup, n int, tcp bool, maxg chan struct{}) {
	defer wg.Done()
	<-maxg
	//time.Sleep(4 * time.Second) //测试协程
	h.Ping(n, tcp)
	fmt.Printf("%v %v\n", h.Ipaddr, h.Reachable)
	maxg <- struct{}{}
}
