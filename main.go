package main

import (
	"bufio"
	"fmt"
	"github.com/go-ping/ping"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

type host struct {
	ipaddr    string
	reachable bool
}

func New(ipaddr string) *host {
	return &host{
		ipaddr: ipaddr,
	}
}

func Ping(h *host) {
	pinger, err := ping.NewPinger(h.ipaddr)
	if err != nil {
		log.Println(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Duration(time.Second)

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
		h.reachable = true
	}
}
func main() {
	var wg sync.WaitGroup
	ipFilePath := "ip.txt"
	ipList := getIPList(ipFilePath)
	var hostList []*host

	for _, ip := range *ipList {
		h := New(ip)
		hostList = append(hostList, h)
		go func() {
			defer wg.Done()
			Ping(h)
			fmt.Printf("%v:%v\n", h.ipaddr, h.reachable)
		}()
		wg.Add(1)
	}
	wg.Wait()
	for _, h := range hostList {
		fmt.Printf("%v:%v\n", h.ipaddr, h.reachable)
	}

}

func getIPList(filepath string) *[]string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reg := regexp.MustCompile(`^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`)
	var ipList []string
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		strLine := string(line)
		if reg.MatchString(strLine) {
			ipList = append(ipList, strLine)
		}
	}
	return &ipList
}
