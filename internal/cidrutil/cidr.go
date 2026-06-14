package cidrutil

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/SiYue-ZO/netkit/pkg/types"
)

// ParseCIDR 解析 CIDR 并返回网段信息
func ParseCIDR(cidrStr string) (*types.CIDRInfo, error) {
	// 如果没有 /前缀，尝试推断
	if !strings.Contains(cidrStr, "/") {
		ip := net.ParseIP(cidrStr)
		if ip == nil {
			return nil, fmt.Errorf("无效的 IP 或 CIDR: %s", cidrStr)
		}
		if ip.To4() != nil {
			cidrStr += "/32"
		} else {
			cidrStr += "/128"
		}
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("无效的 CIDR: %s: %w", cidrStr, err)
	}

	ipv4 := ipNet.IP.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("暂不支持 IPv6 CIDR: %s", cidrStr)
	}

	ones, bits := ipNet.Mask.Size()
	hostCount := int(math.Pow(2, float64(bits-ones))) - 2 // 减去网络地址和广播地址
	if ones >= 31 {
		hostCount = int(math.Pow(2, float64(bits-ones)))
	}

	firstIP := make(net.IP, 4)
	copy(firstIP, ipv4)
	// 第一个可用 IP = 网络地址 + 1 (对于 /30 及更大的网段)
	if ones < 31 {
		firstIP[3]++
	}

	lastIP := make(net.IP, 4)
	binary.BigEndian.PutUint32(lastIP, binary.BigEndian.Uint32(ipv4)+uint32(math.Pow(2, float64(bits-ones)))-1)

	broadcast := ""
	if ones < 31 {
		broadcast = lastIP.String()
	}

	return &types.CIDRInfo{
		CIDR:      cidrStr,
		Network:   ipNet.IP.String(),
		FirstIP:   firstIP.String(),
		LastIP:    lastIP.String(),
		Mask:      fmt.Sprintf("%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3]),
		HostCount: hostCount,
		Broadcast: broadcast,
	}, nil
}

// ExpandCIDR 展开 CIDR 为 IP 列表
func ExpandCIDR(cidrStr string) ([]string, error) {
	info, err := ParseCIDR(cidrStr)
	if err != nil {
		return nil, err
	}

	if info.HostCount > 65536 {
		return nil, fmt.Errorf("CIDR 范围过大 (%d 个主机)，最多支持 65536", info.HostCount)
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, err
	}

	var ips []string
	ip := make(net.IP, 4)
	copy(ip, ipNet.IP.To4())

	start := binary.BigEndian.Uint32(ip)
	ones, bits := ipNet.Mask.Size()
	end := start + uint32(math.Pow(2, float64(bits-ones)))

	// 跳过网络地址和广播地址
	if ones < 31 {
		start++
		end--
	}

	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}

	return ips, nil
}

// AggregateCIDRs 将多个 CIDR 聚合为超网
func AggregateCIDRs(cidrs []string) ([]string, error) {
	var networks []*net.IPNet
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("无效的 CIDR %s: %w", cidr, err)
		}
		networks = append(networks, ipNet)
	}

	// 简单聚合：尝试合并相邻网段
	result := make([]string, 0, len(networks))
	merged := make(map[string]bool)

	for i, n1 := range networks {
		if merged[n1.String()] {
			continue
		}
		for j, n2 := range networks {
			if i == j || merged[n2.String()] {
				continue
			}
			// 检查 n2 是否是 n1 的子网
			if n1.Contains(n2.IP) {
				merged[n2.String()] = true
			}
		}
		result = append(result, n1.String())
	}

	return result, nil
}
