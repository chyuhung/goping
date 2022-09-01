package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/chyuhung/goping/Host"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var hostList []*Host.Host
	//输入ip文件
	var input = flag.String("i", "ip.txt", "ip输入文件")
	//输出结果文件
	var output = flag.String("o", "result.txt", "结果输出文件")
	//最大协程数量
	var maxNum = flag.Int("m", 10, "协程数,不超过ip数量,超过无效")
	//ping包数量
	var pingNum = flag.Int("n", 1, "ping包数")
	flag.Parse()
	fmt.Println("协程数:", *maxNum)
	fmt.Println("ping包数:", *pingNum)
	fmt.Println("读取文件:", *input)
	fmt.Println("输出文件:", *output)

	ipList := getIPList(*input)
	//限制最大协程数
	maxg := make(chan struct{}, *maxNum)
	count := len(*ipList) //获取ip数量
	for i := 0; i < *maxNum; i++ {
		if count <= 0 { //最大不超过ip数量
			break
		}
		maxg <- struct{}{}
	}
	for _, ip := range *ipList {
		h := Host.NewHost(ip)
		hostList = append(hostList, h)
		wg.Add(1)
		go h.Run(&wg, *pingNum, maxg)
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
