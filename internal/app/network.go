package app

import (
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
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				normalized := net.ParseIP(addr.String()).String()
				if isIPv6Unique(normalized, &ipv6) {
					ipv6 = append(ipv6, normalized)
				}
			} else {
				normalized := ip.String()
				if isIPv4Unique(normalized, &ipv4) {
					ipv4 = append(ipv4, normalized)
				}
			}
		}
	}
	sort.Strings(ipv4)
	sort.Strings(ipv6)
	return
}

func isIPv4Unique(ip string, existing *[]string) bool {
	for _, v := range *existing {
		if v == ip {
			return false
		}
	}
	return true
}

func isIPv6Unique(ip string, existing *[]string) bool {
	for _, v := range *existing {
		if net.ParseIP(v).Equal(net.ParseIP(ip)) {
			return false
		}
	}
	return true
}

func isVirtualInterface(name string) bool {
	prefixes := []string{"docker", "veth", "br-", "virbr", "vnet", "tap", "tun", "vpn", "wg"}
	name = strings.ToLower(name)
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func LocalIPv4() string {
	ipvs, _ := LocalIPs()
	if len(ipvs) == 0 {
		return ""
	}
	return strings.Join(ipvs, ",")
}

func LocalIPv6() string {
	_, ipvs := LocalIPs()
	if len(ipvs) == 0 {
		return ""
	}
	return strings.Join(ipvs, ",")
}

// PublicIPs 获取公网 IPv4 和 IPv6 地址
func PublicIPs() (ipv4, ipv6 string) {
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

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			continue
		}

		isIPv4 := parsedIP.To4() != nil
		isIPv6 := !isIPv4

		if svc.ipv6 {
			if !isIPv6 {
				continue
			}
			if ipv6 == "" {
				ipv6 = ip
			}
		} else {
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
