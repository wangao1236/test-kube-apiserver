package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

const ConcurrentNum = 50

var client *http.Client

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout     int = 90
	MaxConnsPerHost     int = 20
)

func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	fmt.Println("connect done, use", conn.LocalAddr().String())

	return conn, err
}

func doGet(client *http.Client, url string, id int) {
	for i := 0; i < 500; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		buf, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("%d::%d: %s -- %v\n", id, i, string(buf), err)
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}
}

func doPost(client *http.Client, url string, id int) {
	for i := 0; i < 500; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		buf, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("%d::%d: %s -- %v\n", id, i, string(buf), err)
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}
}

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			//Proxy: http.ProxyFromEnvironment,
			//DialContext: (&net.Dialer{
			//	Timeout:   30 * time.Second,
			//	KeepAlive: 30 * time.Second,
			//}).DialContext,
			MaxConnsPerHost:     MaxConnsPerHost,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
			Dial:                PrintLocalDial,
		},
	}
	//client = &http.Client{
	//	Transport: &http.Transport{
	//		Dial: PrintLocalDial,
	//	},
	//}
}

//提供给多协程调用
func Fetch(dstUrl string, method string) {
	//resp, _ := client.Get(dstUrl)
	//buf, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("resp: %+v -- %+v\n", string(buf), err)
	//if err := resp.Body.Close(); err != nil {
	//	fmt.Println(err)
	//}

	req, _ := http.NewRequest(method, dstUrl, nil)
	resp, _ := client.Do(req)
	//defer func() {
	//	err := resp.Body.Close()
	//	fmt.Printf("resp close err: %+v\n", err)
	//}()

	fmt.Printf("resp: %+v\n", resp.StatusCode)
	_, _ = ioutil.ReadAll(resp.Body)
	//r, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(resp.Close)
	//fmt.Printf("resp: %+v\n", resp.Status)
	//if err := resp.Body.Close(); err != nil {
	//	fmt.Println(err)
	//}
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(ConcurrentNum)
	for i := 0; i < ConcurrentNum; i++ {
		go func(idx int, wg *sync.WaitGroup) {
			for j := 0; j < 50; j++ {
				Fetch("http://127.0.0.1:12345/post", "POST")
				fmt.Printf(">>>>Transport: %p\n", &client.Transport)
			}
			wg.Done()
		}(i, wg)
	}

	//for i := 0 ; i < ConcurrentNum ; i++ {
	//	go func() {
	//		doGet(client, "http://lb2:12345/get", 1)
	//		wg.Done()
	//	}()
	//}
	wg.Wait()
	//time.Sleep(10 * time.Second)
}
