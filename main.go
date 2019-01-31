package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func resolveAddress(host string, port int, timeout time.Duration) (string, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", host, port), timeout)
	if err == nil {
		conn.Close()
		return host, nil
	}
	_, addr, err := net.LookupSRV("", "", host)
	if err != nil {
		return host, fmt.Errorf("Failed to lookup '%v'. Error: %v.", host, err)
	}
	if len(addr) <= 0 {
		return host, fmt.Errorf("Failed to lookup '%v'. No SRV records returned. Error: %v.", host, err)
	}
	for i := range addr {
		if len(addr[i].Target) <= 0 {
			continue
		}
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", addr[i].Target, port), timeout)
		if err == nil {
			conn.Close()
			return addr[i].Target, nil
		}
	}
	return host, fmt.Errorf("Failed to resolve '%v'. No alternatives found.", host)
}

func main() {
	var addrToResolve string
	var port int
	var timeout int
	flag.StringVar(&addrToResolve, "addr", "", "Address to resolve")
	flag.IntVar(&port, "port", 0, "Port to connect")
	flag.IntVar(&timeout, "timeoutsec", 5, "Timeout seconds for connection attempt")
	flag.Parse()
	if len(addrToResolve) <= 0 {
		log.Printf("Empty address specified.")
		os.Exit(1)
	}
	if port <= 0 {
		log.Printf("Invalid port specified.")
		os.Exit(1)
	}
	if timeout <= 0 {
		log.Printf("Invalid timeout specified.")
		os.Exit(1)
	}
	result, err := resolveAddress(addrToResolve, port, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Printf("Failed to resolve address '%v' with port '%v'. Error: %v", addrToResolve, port, err)
		os.Exit(1)
	}
	fmt.Printf(result)
	os.Exit(0)
}
