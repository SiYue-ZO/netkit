package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SiYue-ZO/netkit/pkg/pool"
	"github.com/SiYue-ZO/netkit/pkg/types"
)

// Top100Ports 常见 Top 100 端口
var Top100Ports = []int{
	7, 9, 13, 21, 22, 23, 25, 26, 37, 53, 79, 80, 81, 88, 106, 110, 111, 113,
	119, 135, 139, 143, 144, 179, 199, 389, 427, 443, 444, 445, 465, 513, 514,
	515, 543, 544, 548, 554, 587, 631, 646, 873, 990, 993, 995, 1025, 1026,
	1027, 1028, 1029, 1110, 1433, 1720, 1723, 1755, 1900, 2000, 2001, 2049,
	2121, 2717, 3000, 3128, 3306, 3389, 3986, 4899, 5000, 5009, 5051, 5060,
	5101, 5190, 5357, 5432, 5631, 5666, 5800, 5900, 6000, 6001, 6646, 7070,
	8000, 8008, 8009, 8080, 8081, 8443, 8888, 9100, 9999, 10000, 27017, 32768,
	49152, 49153, 49154, 49155, 49156, 49157,
}

// Top1000Ports 常见 Top 1000 端口 (省略，使用 Top100 + 常用端口)
// 完整列表可在后续版本补充

// Scanner 端口扫描器
type Scanner struct {
	Timeout  time.Duration
	Threads  int
	Ports    []int
	Protocol string // "tcp" or "udp"
	Verbose  bool
}

// NewScanner 创建扫描器
func NewScanner(timeout time.Duration, threads int, protocol string, verbose bool) *Scanner {
	return &Scanner{
		Timeout:  timeout,
		Threads:  threads,
		Protocol: protocol,
		Verbose:  verbose,
	}
}

// Scan 扫描单个主机的指定端口
func (s *Scanner) Scan(host string) *types.HostResult {
	result := &types.HostResult{
		Host: host,
	}

	// 解析 IP
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupHost(host)
		if err != nil || len(ips) == 0 {
			return result
		}
		result.IP = ips[0]
	} else {
		result.IP = host
	}

	var mu sync.Mutex
	wp := pool.New(s.Threads)

	for _, port := range s.Ports {
		p := port
		wp.Submit(func() {
			addr := net.JoinHostPort(host, strconv.Itoa(p))
			conn, err := net.DialTimeout(s.Protocol, addr, s.Timeout)
			if err == nil {
				conn.Close()
				mu.Lock()
				result.Ports = append(result.Ports, types.PortInfo{
					Port:     p,
					Protocol: s.Protocol,
					State:    "open",
					Service:  WellKnownPort(p),
				})
				mu.Unlock()
			}
		})
	}

	wp.Wait()

	if len(result.Ports) > 0 {
		result.Alive = true
	}

	return result
}

// ScanMultiple 批量扫描
func (s *Scanner) ScanMultiple(targets []string) []types.HostResult {
	var results []types.HostResult
	var mu sync.Mutex

	wp := pool.New(s.Threads)

	for _, target := range targets {
		t := target
		wp.Submit(func() {
			r := s.Scan(t)
			mu.Lock()
			results = append(results, *r)
			mu.Unlock()
		})
	}

	wp.Wait()
	return results
}

// ParsePorts 解析端口字符串，支持: 80,443,1-1000,top100,top1000,full
func ParsePorts(portStr string) ([]int, error) {
	if portStr == "" {
		return Top100Ports, nil
	}

	switch portStr {
	case "top100":
		return Top100Ports, nil
	case "top1000":
		return append(Top100Ports, expandTop1000()...), nil
	case "full", "all", "-":
		return parseRange("1-65535")
	}

	var ports []int
	seen := make(map[int]bool)

	parts := strings.Split(portStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangePorts, err := parseRange(part)
			if err != nil {
				return nil, err
			}
			for _, p := range rangePorts {
				if !seen[p] {
					seen[p] = true
					ports = append(ports, p)
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("无效端口号: %s", part)
			}
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("端口号超出范围: %d", p)
			}
			if !seen[p] {
				seen[p] = true
				ports = append(ports, p)
			}
		}
	}

	return ports, nil
}

func parseRange(r string) ([]int, error) {
	parts := strings.Split(r, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效端口范围: %s", r)
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("无效起始端口: %s", parts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("无效结束端口: %s", parts[1])
	}

	if start > end {
		return nil, fmt.Errorf("起始端口大于结束端口: %d > %d", start, end)
	}

	if start < 1 || end > 65535 {
		return nil, fmt.Errorf("端口范围超出有效范围: %d-%d", start, end)
	}

	ports := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return ports, nil
}

// WellKnownPort 返回常见端口对应的服务名
func WellKnownPort(port int) string {
	services := map[int]string{
		21: "ftp", 22: "ssh", 23: "telnet", 25: "smtp", 53: "dns",
		80: "http", 110: "pop3", 111: "rpcbind", 135: "msrpc",
		139: "netbios-ssn", 143: "imap", 443: "https", 445: "microsoft-ds",
		993: "imaps", 995: "pop3s", 1433: "mssql", 1521: "oracle",
		3306: "mysql", 3389: "ms-wbt-server", 5432: "postgresql",
		5900: "vnc", 6379: "redis", 8080: "http-proxy", 8443: "https-alt",
		27017: "mongodb",
	}
	if svc, ok := services[port]; ok {
		return svc
	}
	return ""
}

// expandTop1000 补充 Top1000 中不在 Top100 的端口
func expandTop1000() []int {
	return []int{
		19, 69, 102, 107, 109, 1434, 1521, 220, 363, 401, 425, 465, 512,
		524, 530, 532, 540, 555, 593, 636, 639, 777, 783, 808, 843, 880,
		981, 987, 1010, 1035, 1040, 1048, 1050, 1053, 1054, 1056, 1058,
		1064, 1065, 1071, 1080, 1094, 1098, 1099, 1100, 1101, 1104, 1106,
		1108, 1111, 1112, 1114, 1117, 1119, 1121, 1122, 1123, 1124, 1126,
		1130, 1131, 1132, 1137, 1138, 1141, 1145, 1147, 1148, 1149, 1151,
		1152, 1154, 1163, 1164, 1165, 1166, 1169, 1174, 1175, 1183, 1185,
		1186, 1187, 1192, 1198, 1199, 1201, 1213, 1216, 1217, 1233, 1234,
		1236, 1244, 1247, 1248, 1259, 1271, 1272, 1277, 1287, 1300, 1301,
		1309, 1310, 1311, 1322, 1328, 1334, 1352, 1417, 1433, 1434, 1443,
		1455, 1461, 1494, 1500, 1501, 1503, 1521, 1524, 1525, 1526, 1527,
		1528, 1529, 1530, 1531, 1532, 1533, 1534, 1535, 1536, 1537, 1538,
		1539, 1540, 1541, 1542, 1543, 1544, 1545, 1546, 1547, 1548, 1549,
		1550, 1551, 1552, 1553, 1554, 1555, 1556, 1557, 1558, 1559, 1560,
	}
}
