package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/chyuhung/goping/Host"
)

func isTrue(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func main() {
	var wg sync.WaitGroup
	var hostList []*Host.Host
	//输入ip文件
	var input = flag.String("i", "ip.txt", "Input file.")
	//输出结果文件
	var output = flag.String("o", "result.txt", "Output file.")
	//最大协程数量
	var maxNum = flag.Int("m", 10, "The number of goroutines, the maximum is the number of IPs.")
	//ping包数量
	var pingNum = flag.Int("n", 1, "Number of messages per ping.")
	var isICMP = flag.Int("t", 1, "Set the type of ping pinger will send.\n"+
		"False means pinger will send an \"unprivileged\" UDP ping.\n"+
		"True means pinger will send a \"privileged\" raw ICMP ping.\n"+
		"NOTE: setting to true requires that it be run with super-user privileges.")
	flag.Parse()

	fmt.Println("Goroutines:", *maxNum)
	fmt.Println("Messages:", *pingNum)
	fmt.Println("Input file:", *input)
	fmt.Println("Output file:", *output)
	fmt.Println("Send privileged ping:", *isICMP, "\n")

	var maxg chan struct{} //限制最大协程数
	ipList := getIPList(*input)
	count := len(*ipList) //获取ip数量
	if count <= *maxNum {
		maxg = make(chan struct{}, count)
	} else {
		maxg = make(chan struct{}, *maxNum)
		count = *maxNum
	}
	for i := count; i > 0; i-- {
		maxg <- struct{}{}
	}
	for _, ip := range *ipList {
		h := Host.NewHost(ip)
		hostList = append(hostList, h)
		wg.Add(1)
		go h.Run(&wg, *pingNum, isTrue(*isICMP), maxg)
	}
	wg.Wait()
	f, err := os.Create(*output)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	for i := range hostList {
		_, err := f.WriteString(hostList[i].Ipaddr + " " + hostList[i].Reachable + "\n")
		if err != nil {
			log.Fatal(err)
		}
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
