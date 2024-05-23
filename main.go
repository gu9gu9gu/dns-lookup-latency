package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <domain> <interval-seconds> <record-type>\n", os.Args[0])
		fmt.Println("Record type can be 'A' or 'AAAA'")
		os.Exit(1)
	}

	domain := os.Args[1]
	intervalStr := os.Args[2]
	recordType := os.Args[3]

	interval, err := strconv.ParseFloat(intervalStr, 64)
	if err != nil {
		fmt.Printf("failed to parse interval(%s): %v\n", intervalStr, err)
		os.Exit(1)
	}

	duration := time.Millisecond * time.Duration(interval*1000)

	for {
		resolveDomain(domain, recordType)
		time.Sleep(duration)
	}
}

func resolveDomain(domain, recordType string) {
	startTime := time.Now()
	var ips []net.IP
	var err error

	switch recordType {
	case "A":
		ips, err = net.LookupIP(domain)
	case "AAAA":
		ips, err = net.LookupAAAA(domain)
	default:
		fmt.Printf("Invalid record type: %s. Use 'A' or 'AAAA'.\n", recordType)
		return
	}

	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Printf("%v : %v\n", startTime.Format(time.StampMilli), err)
		return
	}

	fmt.Printf("%v : %.3f ms\n", startTime.Format(time.StampMilli), elapsed.Seconds()*1000)
	for _, ip := range ips {
		fmt.Println(ip)
	}
}