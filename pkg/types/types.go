package types

// HostResult 单个主机的扫描结果
type HostResult struct {
	Host    string      `json:"host"`
	IP      string      `json:"ip,omitempty"`
	Alive   bool        `json:"alive,omitempty"`
	RTT     string      `json:"rtt,omitempty"`
	Ports   []PortInfo  `json:"ports,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// PortInfo 端口信息
type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"` // tcp/udp
	State    string `json:"state"`    // open/closed/filtered
	Service  string `json:"service,omitempty"`
	Banner   string `json:"banner,omitempty"`
}

// ScanOptions 通用扫描选项
type ScanOptions struct {
	Targets  []string
	Timeout  int // seconds
	Threads  int
	Verbose  bool
	Output   string // output file path
	Format   string // json/csv/table
	NoColor  bool
}

// CIDRInfo CIDR 网段信息
type CIDRInfo struct {
	CIDR         string   `json:"cidr"`
	Network      string   `json:"network"`
	FirstIP      string   `json:"first_ip"`
	LastIP       string   `json:"last_ip"`
	Mask         string   `json:"mask"`
	HostCount    int      `json:"host_count"`
	Broadcast    string   `json:"broadcast,omitempty"`
	IPs          []string `json:"ips,omitempty"`
}
