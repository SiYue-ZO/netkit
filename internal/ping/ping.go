package ping

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/SiYue-ZO/netkit/pkg/pool"
	"github.com/SiYue-ZO/netkit/pkg/types"
)

// Pinger 主机存活探测器
type Pinger struct {
	Timeout time.Duration
	Count   int
	Method  string // "icmp" or "tcp"
}

// NewPinger 创建探测器
func NewPinger(timeout time.Duration, count int, method string) *Pinger {
	return &Pinger{
		Timeout: timeout,
		Count:   count,
		Method:  method,
	}
}

// Ping 探测单个主机
func (p *Pinger) Ping(host string) *types.HostResult {
	result := &types.HostResult{
		Host: host,
	}

	// 解析 IP
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupHost(host)
		if err != nil || len(ips) == 0 {
			result.Alive = false
			return result
		}
		ip = net.ParseIP(ips[0])
		result.IP = ips[0]
	} else {
		result.IP = host
	}

	switch p.Method {
	case "icmp":
		p.pingICMP(ip, result)
	case "tcp":
		p.pingTCP(ip, result)
	default:
		// 默认先尝试 ICMP，失败后尝试 TCP
		p.pingICMP(ip, result)
		if !result.Alive {
			p.pingTCP(ip, result)
		}
	}

	return result
}

// PingMultiple 批量探测
func (p *Pinger) PingMultiple(targets []string, threads int) []types.HostResult {
	var results []types.HostResult
	var mu sync.Mutex

	wp := pool.New(threads)

	for _, target := range targets {
		t := target
		wp.Submit(func() {
			r := p.Ping(t)
			mu.Lock()
			results = append(results, *r)
			mu.Unlock()
		})
	}

	wp.Wait()
	return results
}

func (p *Pinger) pingICMP(ip net.IP, result *types.HostResult) {
	// Windows 和 Linux 下 ICMP 需要特权，这里使用简化实现
	// 尝试发送 ICMP Echo Request
	if runtime.GOOS == "windows" {
		// Windows 下 ICMP 需要管理员权限，回退到 TCP
		p.pingTCP(ip, result)
		return
	}

	// Unix 下尝试 raw socket ICMP
	conn, err := net.DialTimeout("ip4:icmp", ip.String(), p.Timeout)
	if err != nil {
		result.Alive = false
		return
	}
	defer conn.Close()

	start := time.Now()
	_ = conn.SetDeadline(time.Now().Add(p.Timeout))

	// 发送 ICMP Echo (type 8, code 0)
	msg := make([]byte, 8)
	msg[0] = 8 // Type: Echo Request
	msg[1] = 0 // Code
	msg[2] = 0 // Checksum (placeholder)
	msg[3] = 0
	msg[4] = 0 // Identifier
	msg[5] = 1
	msg[6] = 0 // Sequence
	msg[7] = 1

	// 计算校验和
	msg[2], msg[3] = checksum(msg)

	_, err = conn.Write(msg)
	if err != nil {
		result.Alive = false
		return
	}

	reply := make([]byte, 20+8)
	_, err = conn.Read(reply)
	if err != nil {
		result.Alive = false
		return
	}

	result.Alive = true
	result.RTT = time.Since(start).Round(time.Millisecond).String()
}

func checksum(msg []byte) (hi, lo byte) {
	sum := 0
	for i := 0; i < len(msg)-1; i += 2 {
		sum += int(msg[i])<<8 | int(msg[i+1])
	}
	if len(msg)%2 == 1 {
		sum += int(msg[len(msg)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	sum = ^sum

	return byte(sum >> 8), byte(sum & 0xff)
}

func (p *Pinger) pingTCP(ip net.IP, result *types.HostResult) {
	// 尝试常见端口
	ports := []int{80, 443, 22, 8080}

	for _, port := range ports {
		addr := net.JoinHostPort(ip.String(), strconv.Itoa(port))
		start := time.Now()
		conn, err := net.DialTimeout("tcp", addr, p.Timeout)
		if err == nil {
			conn.Close()
			result.Alive = true
			result.RTT = time.Since(start).Round(time.Millisecond).String()
			return
		}
	}

	result.Alive = false
}

// ResolveTargets 解析目标列表（支持 CIDR/文件/单 IP）
func ResolveTargets(targets []string, inputFile string) ([]string, error) {
	var allTargets []string

	// 从参数获取
	for _, t := range targets {
		// 检查是否是 CIDR
		if isCIDR(t) {
			ips, err := expandCIDRSimple(t)
			if err != nil {
				return nil, fmt.Errorf("展开 CIDR %s 失败: %w", t, err)
			}
			allTargets = append(allTargets, ips...)
		} else {
			allTargets = append(allTargets, t)
		}
	}

	// 从文件读取
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("读取文件 %s 失败: %w", inputFile, err)
		}
		lines := splitLines(string(data))
		for _, line := range lines {
			line = trimSpace(line)
			if line == "" || line[0] == '#' {
				continue
			}
			if isCIDR(line) {
				ips, err := expandCIDRSimple(line)
				if err != nil {
					return nil, err
				}
				allTargets = append(allTargets, ips...)
			} else {
				allTargets = append(allTargets, line)
			}
		}
	}

	return allTargets, nil
}

func isCIDR(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

func expandCIDRSimple(cidr string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		ips = append(ips, ip.String())
	}
	// 移除网络地址和广播地址
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}
	return ips, nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if line != "" && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
