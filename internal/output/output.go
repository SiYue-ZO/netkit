package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/SiYue-ZO/netkit/pkg/types"
	"github.com/fatih/color"
)

// Writer 统一输出接口
type Writer struct {
	format  string
	noColor bool
	file    *os.File
	writer  io.Writer
}

// NewWriter 创建输出写入器
func NewWriter(format, outputPath string, noColor bool) (*Writer, error) {
	w := &Writer{
		format:  format,
		noColor: noColor,
		writer:  os.Stdout,
	}

	if noColor {
		color.NoColor = true
	}

	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("创建输出文件失败: %w", err)
		}
		w.file = f
		w.writer = f
	}

	return w, nil
}

// Write 实现 io.Writer 接口
func (w *Writer) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

// Close 关闭输出
func (w *Writer) Close() {
	if w.file != nil {
		w.file.Close()
	}
}

// WriteHostResults 写入主机结果
func (w *Writer) WriteHostResults(results []types.HostResult) error {
	switch w.format {
	case "json":
		return w.writeJSON(results)
	case "csv":
		return w.writeCSV(results)
	default:
		return w.writeTable(results)
	}
}

// WriteCIDRInfo 写入 CIDR 信息
func (w *Writer) WriteCIDRInfo(info *types.CIDRInfo) error {
	switch w.format {
	case "json":
		return w.writeJSON(info)
	case "csv":
		return w.writeCIDRCSV(info)
	default:
		return w.writeCIDRTable(info)
	}
}

func (w *Writer) writeJSON(data interface{}) error {
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (w *Writer) writeCSV(results []types.HostResult) error {
	csvWriter := csv.NewWriter(w.writer)
	defer csvWriter.Flush()

	// 写入表头
	header := []string{"Host", "IP", "Alive", "RTT", "Port", "Protocol", "State", "Service", "Banner"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, r := range results {
		if len(r.Ports) == 0 {
			row := []string{r.Host, r.IP, fmt.Sprintf("%v", r.Alive), r.RTT, "", "", "", "", ""}
			if err := csvWriter.Write(row); err != nil {
				return err
			}
		}
		for _, p := range r.Ports {
			row := []string{
				r.Host, r.IP, fmt.Sprintf("%v", r.Alive), r.RTT,
				fmt.Sprintf("%d", p.Port), p.Protocol, p.State, p.Service, p.Banner,
			}
			if err := csvWriter.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Writer) writeTable(results []types.HostResult) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	for _, r := range results {
		// 主机行
		status := red("down")
		if r.Alive {
			status = green("up")
		}

		fmt.Fprintf(w.writer, "%s %s %s", cyan(r.Host), status, r.RTT)
		if r.IP != "" && r.IP != r.Host {
			fmt.Fprintf(w.writer, " (%s)", r.IP)
		}
		fmt.Fprintln(w.writer)

		// 端口行
		for _, p := range r.Ports {
			stateStr := red("closed")
			if p.State == "open" {
				stateStr = green("open")
			}
			line := fmt.Sprintf("  %-6d %-6s %s", p.Port, p.Protocol, stateStr)
			if p.Service != "" {
				line += fmt.Sprintf(" %s", p.Service)
			}
			if p.Banner != "" {
				banner := p.Banner
				if len(banner) > 60 {
					banner = banner[:60] + "..."
				}
				banner = strings.ReplaceAll(banner, "\n", " ")
				banner = strings.ReplaceAll(banner, "\r", "")
				line += fmt.Sprintf(" [%s]", banner)
			}
			fmt.Fprintln(w.writer, line)
		}
	}

	return nil
}

func (w *Writer) writeCIDRTable(info *types.CIDRInfo) error {
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Fprintf(w.writer, "CIDR:       %s\n", cyan(info.CIDR))
	fmt.Fprintf(w.writer, "网络地址:   %s\n", info.Network)
	fmt.Fprintf(w.writer, "掩码:       %s\n", info.Mask)
	fmt.Fprintf(w.writer, "起始 IP:    %s\n", info.FirstIP)
	fmt.Fprintf(w.writer, "结束 IP:    %s\n", info.LastIP)
	fmt.Fprintf(w.writer, "主机数量:   %d\n", info.HostCount)
	if info.Broadcast != "" {
		fmt.Fprintf(w.writer, "广播地址:   %s\n", info.Broadcast)
	}
	return nil
}

func (w *Writer) writeCIDRCSV(info *types.CIDRInfo) error {
	csvWriter := csv.NewWriter(w.writer)
	defer csvWriter.Flush()

	header := []string{"CIDR", "Network", "Mask", "FirstIP", "LastIP", "HostCount", "Broadcast"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	row := []string{info.CIDR, info.Network, info.Mask, info.FirstIP, info.LastIP,
		fmt.Sprintf("%d", info.HostCount), info.Broadcast}
	return csvWriter.Write(row)
}

// Info 输出信息行
func (w *Writer) Info(format string, args ...interface{}) {
	if !w.noColor {
		color.Cyan(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

// Warn 输出警告行
func (w *Writer) Warn(format string, args ...interface{}) {
	if !w.noColor {
		color.Yellow(format, args...)
		return
	}
	fmt.Printf(format+"\n", args...)
}

// Error 输出错误行
func (w *Writer) Error(format string, args ...interface{}) {
	if !w.noColor {
		color.Red(format, args...)
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
