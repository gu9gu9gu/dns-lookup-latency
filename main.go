package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// resolveDomain 根据提供的记录类型（A或AAAA）查询域名并打印IP地址
func resolveDomain(ctx context.Context, domain, recordType string) {
	startTime := time.Now()
	var ips []net.IP
	var err error

	switch recordType {
	case "AAAA":
		ips, err = net.LookupIP(domain)
		if err != nil {
			fmt.Printf("%v : %v\n", startTime.Format("Jan _2 15:04:05"), err)
			return
		}
		ips = filterIPv6(ips)
		printIPs(startTime, ips)
	case "A":
		ips, err = net.LookupIP(domain)
		if err != nil {
			fmt.Printf("%v : %v\n", startTime.Format("Jan _2 15:04:05"), err)
			return
		}
		ips = filterIPv4(ips)
		printIPs(startTime, ips)
	default:
		// 查询所有记录
		ips, err = net.LookupIP(domain)
		if err != nil {
			fmt.Printf("%v : %v\n", startTime.Format("Jan _2 15:04:05"), err)
			return
		}
		printIPs(startTime, ips)
	}
}

// filterIPv4 从IP列表中过滤出IPv4地址
func filterIPv4(ips []net.IP) []net.IP {
	ipv4s := []net.IP{}
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4s = append(ipv4s, ip)
		}
	}
	return ipv4s
}

// filterIPv6 从IP列表中过滤出IPv6地址
func filterIPv6(ips []net.IP) []net.IP {
	ipv6s := []net.IP{}
	for _, ip := range ips {
		if ip.To4() == nil { // To4() 对IPv6地址返回nil
			ipv6s = append(ipv6s, ip)
		}
	}
	return ipv6s
}

// printIPs 打印IP地址及查询时间
func printIPs(startTime time.Time, ips []net.IP) {
	elapsed := time.Since(startTime)
	fmt.Printf("%v : %.3f ms\n", startTime.Format("Jan _2 15:04:05"), elapsed.Seconds()*1000)
	for _, ip := range ips {
		fmt.Println(ip)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <domain> <interval-seconds> [<record-type>]\n", os.Args[0])
		fmt.Println("Record type can be 'A', 'AAAA', or omit for both.")
		os.Exit(1)
	}

	domain := os.Args[1]
	intervalStr := os.Args[2]

	interval, err := strconv.ParseFloat(intervalStr, 64)
	if err != nil {
		fmt.Printf("Failed to parse interval: %v\n", err)
		os.Exit(1)
	}

	recordType := "" // 默认查询所有记录
	if len(os.Args) > 3 {
		recordType = os.Args[3]
		if recordType != "A" && recordType != "AAAA" {
			fmt.Printf("Invalid record type: %s. Use 'A' or 'AAAA'.\n", recordType)
			os.Exit(1)
		}
	}

	// 创建一个ticker，用于周期性执行
	ticker := time.NewTicker(time.Duration(interval*1000) * time.Millisecond)
	defer ticker.Stop() // 确保在程序结束时停止ticker

	// 使用context来管理超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting after 10 seconds.")
			return
		case <-ticker.C:
			resolveDomain(ctx, domain, recordType)
		}
	}
}