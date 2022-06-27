package main

import (
	"bufio"
	"github.com/go-ping/ping"
	"io"
	"log"
	"os"
	"regexp"
)

func main() {
	ipfile := "ip.txt"
	ipList := getIPList(ipfile)
	for _, ip := range ipList {
		pinger, err := ping.NewPinger(ip)
		if err != nil {
			log.Fatal(err)
		}
		pinger.Size = 24
		pinger.Count = 2
		go func() {
			err := pinger.Run()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

}

type host struct {
	ipaddr    string
	reachable bool
}

func getIPList(filepath string) []string {
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
	return ipList
}
