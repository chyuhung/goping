package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func main() {
	ipfile := "ip.txt"
	ipList := getIPListFromFile(ipfile)
	for _, ip := range ipList {
		fmt.Println(ip)
	}

}

func getIPListFromFile(filepath string) []string {
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
