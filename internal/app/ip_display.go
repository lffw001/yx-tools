package app

import (
	"fmt"
	"strings"
)

// PrintLocalIPs 在命令行打印本机 IPv4/IPv6 和公网 IP 地址
func PrintLocalIPs() {
	ipv4, ipv6 := LocalIPs()

	fmt.Println("本机 IP 地址")
	fmt.Println(strings.Repeat("-", 40))

	if len(ipv4) > 0 {
		fmt.Printf("IPv4: %s\n", strings.Join(ipv4, ", "))
	} else {
		fmt.Println("IPv4: 未找到")
	}

	if len(ipv6) > 0 {
		fmt.Printf("IPv6: %s\n", strings.Join(ipv6, ", "))
	} else {
		fmt.Println("IPv6: 未找到")
	}

	fmt.Println()
	fmt.Println("公网 IP 地址")
	fmt.Println(strings.Repeat("-", 40))

	pubIPV4, pubIPV6 := PublicIPs()
	if pubIPV4 != "" {
		fmt.Printf("IPv4: %s\n", pubIPV4)
	} else {
		fmt.Println("IPv4: 无法获取")
	}

	if pubIPV6 != "" {
		fmt.Printf("IPv6: %s\n", pubIPV6)
	} else {
		fmt.Println("IPv6: 无法获取")
	}
}
