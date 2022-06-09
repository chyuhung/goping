package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	ipfile := "ip.txt"
	ipList := getIPList(ipfile)
	for _, ip := range ipList {
		func(ip string) {
			fmt.Println(ip)
			cmd := exec.Command("ping", "-c", "1", ip)
			cmd.Start()
		}(ip)
	}

}

func getIPList(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reg := regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
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
