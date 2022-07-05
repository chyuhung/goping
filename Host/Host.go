package Host

import (
	"github.com/go-ping/ping"
	"log"
	"time"
)

type Host struct {
	Ipaddr    string
	Reachable bool
}

func NewHost(ipaddr string) *Host {
	return &Host{
		Ipaddr: ipaddr,
	}
}

func (h *Host) Ping(n int) {
	pinger, err := ping.NewPinger(h.Ipaddr)
	if err != nil {
		log.Println(err)
	}
	pinger.Count = n
	pinger.Timeout = time.Second

	/*
		设置pinger将发送的类型。
		false表示pinger将发送“未经授权”的UDP ping
		true表示pinger将发送“特权”原始ICMP ping
	*/
	pinger.SetPrivileged(true)
	// 运行pinger
	pinger.Run()
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		h.Reachable = true
	}
}
