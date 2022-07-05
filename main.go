package main

import (
	"bufio"
	"gitee.com/chyuhung/goping/GoLimit"
	"gitee.com/chyuhung/goping/Host"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

func main() {
	ipFilePath := "ip.txt"
	resultPath := "result.txt"
	var wg sync.WaitGroup
	var hostList []*Host.Host
	ipList := getIPList(ipFilePath)

	//最大协程数量
	g := GoLimit.NewGoLimit(1)
	//ping包数量
	pingNum := 2

	for _, ip := range *ipList {
		h := Host.NewHost(ip)
		hostList = append(hostList, h)
		goFunc := func() {
			defer wg.Done()
			h.Ping(pingNum)
			log.Printf("ipaddr:%v reachable:%v\n", h.Ipaddr, h.Reachable)
		}
		g.Run(goFunc)
		wg.Add(1)
	}
	wg.Wait()
	f, err := os.Create(resultPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	for _, h := range hostList {
		f.WriteString(h.Ipaddr + ":" + strconv.FormatBool(h.Reachable) + "\n")
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
