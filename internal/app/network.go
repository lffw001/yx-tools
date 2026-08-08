package app

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

// LocalIPs 返回本机所有非回环、非虚拟网卡的有效 IP 地址
func LocalIPs() (ipv4, ipv6 []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range interfaces {
		// 跳过回环、down、虚拟网卡
		if iface.Flags&net.FlagLoopback != 0 ||
			iface.Flags&net.FlagUp == 0 ||
			isVirtualInterface(iface.Name) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			// 过滤回环和私有/链路本地地址（除非明确要求）
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				continue
			}
			// 清理前缀
			ip = ip.To4()
			if ip == nil {
				// IPv6
				normalized := ip.String()
				// 去掉压缩的零段，方便阅读
				if isIPv6Unique(normalized, &ipv6) {
					ipv6 = append(ipv6, normalized)
				}
			} else {
				// IPv4
				normalized := ip.String()
				if isIPv4Unique(normalized, &ipv4) {
					ipv4 = append(ipv4, normalized)
				}
			}
		}
	}
	// 排序让输出更稳定
	sort.Strings(ipv4)
	sort.Strings(ipv6)
	return
}

// isIPv4Unique 去重（同一 IP 可能出现在多个网卡上）
func isIPv4Unique(ip string, existing *[]string) bool {
	for _, v := range *existing {
		if v == ip {
			return false
		}
	}
	return true
}

// isIPv6Unique IPv6 去重，规范化后比较
func isIPv6Unique(ip string, existing *[]string) bool {
	for _, v := range *existing {
		if normalizeIPv6(v) == normalizeIPv6(ip) {
			return false
		}
	}
	return true
}

// normalizeIPv6 将 IPv6 地址规范化（全小写、去掉前导零、压缩零段）
func normalizeIPv6(s string) string {
	ip := net.ParseIP(s)
	if ip == nil {
		return s
	}
	return ip.String()
}

// isVirtualInterface 判断是否为虚拟网卡（Docker、VPN、虚拟机等）
func isVirtualInterface(name string) bool {
	// 常见虚拟网卡前缀
	prefixes := []string{
		"docker", "veth", "br-",       // Docker
		"virbr", "vnet", "tap", "tun", // KVM/QEMU
		"vpn", "wg", "wlan",            // VPN/WireGuard
		"lo",                           // 回环已在上面过滤
	}
	name = lower(name)
	for _, p := range prefixes {
		if len(name) >= len(p) && name[:len(p)] == p {
			return true
		}
	}
	return false
}

func lower(s string) string {
	if s >= "A" && s <= "Z" {
		return s + string(s[0]-'A'+'a')
	}
	return s
}

// LocalIPv4 返回本机所有有效 IPv4 地址（逗号分隔）
func LocalIPv4() string {
	ipvs, _ := LocalIPs()
	if len(ipvs) == 0 {
		return ""
	}
	return joinStrings(ipvs)
}

// LocalIPv6 返回本机所有有效 IPv6 地址（逗号分隔）
func LocalIPv6() string {
	_, ipvs := LocalIPs()
	if len(ipvs) == 0 {
		return ""
	}
	return joinStrings(ipvs)
}

// LocalIPInfo 返回本机的网络信息，供 JSON 序列化
func LocalIPInfo() map[string]string {
	ipv4, ipv6 := LocalIPs()
	return map[string]string{
		"ipv4": joinStrings(ipv4),
		"ipv6": joinStrings(ipv6),
	}
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for i := 1; i < len(ss); i++ {
		out += "," + ss[i]
	}
	return out
}

// PublicIPs 获取公网 IPv4 和 IPv6 地址
func PublicIPs() (ipv4, ipv6 string) {
	// 使用多个服务以提高可靠性
	services := []struct {
		url  string
		ipv6 bool
	}{
		{"https://api.ipify.org", false},
		{"https://ipinfo.io/ip", false},
		{"https://ip.sb", false},
		{"https://checkip.amazonaws.com", false},
		{"https://api64.ipify.org", true},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, svc := range services {
		req, err := http.NewRequest("GET", svc.url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Accept", "text/plain")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		ip := strings.TrimSpace(string(body))
		if ip == "" {
			continue
		}

		// 验证是否为合法的IP地址
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			continue
		}

		// 验证IP类型是否与预期一致
		isIPv4 := parsedIP.To4() != nil
		isIPv6 := !isIPv4

		if svc.ipv6 {
			// 请求IPv6，但返回的是IPv4，跳过
			if !isIPv6 {
				continue
			}
			if ipv6 == "" {
				ipv6 = ip
			}
		} else {
			// 请求IPv4，但返回的是IPv6，跳过
			if !isIPv4 {
				continue
			}
			if ipv4 == "" {
				ipv4 = ip
			}
		}

		if ipv4 != "" && ipv6 != "" {
			break
		}
	}

	return
}

// PublicIPInfo 返回公网IP信息，供 JSON 序列化
func PublicIPInfo() map[string]string {
	ipv4, ipv6 := PublicIPs()
	return map[string]string{
		"ipv4": ipv4,
		"ipv6": ipv6,
	}
}

// PrintLocalIPs 在命令行打印本机和本地 IP 地址
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
