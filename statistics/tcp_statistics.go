package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func ExtractAddrs(log string) (string, string) {
	regex := `^(\d{1,3}\.){3}\d{1,3}:\d{1,5}`
	reg, _ := regexp.Compile(regex)
	addrs := reg.FindAllString(log, 2)
	return addrs[0], addrs[1]
}

func ExtractAddr(log string) string {
	regex := `(\d{1,3}\.){3}\d{1,3}:\d{1,5}`
	reg, _ := regexp.Compile(regex)
	return reg.FindString(log)
}

func main() {

	set := make(map[string]int)

	file, err := os.Open("../log/nginx.log_1")
	//file, err := os.Open("../log/1.68_apiserver.log_1")
	//file, err := os.Open("../log/1.67_apiserver.log_1")
	//file, err := os.Open("../log/test-curl.log_1")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	cnt := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		src := ExtractAddr(string(a))
		set[src]++
		cnt++
	}
	fmt.Printf("tcp ports' num: %d requests' num: %d\n", len(set), cnt)
}
