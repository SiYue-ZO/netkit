package cmd

import (
	"fmt"
	"time"

	"github.com/SiYue-ZO/netkit/internal/output"
	"github.com/SiYue-ZO/netkit/internal/ping"
	"github.com/SiYue-ZO/netkit/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	scanTargets  []string
	scanFile     string
	scanPorts    string
	scanProto    string
	scanSkipPing bool
)

var portscanCmd = &cobra.Command{
	Use:   "portscan",
	Short: "端口扫描 (CONNECT/SYN/UDP)",
	Long: `端口扫描工具，支持:
  - TCP Connect 扫描 (默认，无需特权)
  - TCP SYN 扫描 (需要 root/管理员权限，后续版本)
  - UDP 扫描 (后续版本)
  - Top100/Top1000/全端口/自定义端口
  - 支持 CIDR/文件输入，并发扫描

示例:
  netkit portscan -t 192.168.1.1
  netkit portscan -t 192.168.1.1 -p 22,80,443,8080
  netkit portscan -t 192.168.1.0/24 --top-ports 100
  netkit portscan -t 192.168.1.1 -p 1-65535 -c 200
  netkit portscan -l hosts.txt -p top1000 -o results.json -f json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 解析目标
		targets, err := ping.ResolveTargets(scanTargets, scanFile)
		if err != nil {
			return err
		}
		if len(targets) == 0 {
			return fmt.Errorf("未指定目标，使用 -t 或 -l 指定")
		}

		// 解析端口
		ports, err := scanner.ParsePorts(scanPorts)
		if err != nil {
			return err
		}

		w, err := output.NewWriter(opts.Format, opts.Output, opts.NoColor)
		if err != nil {
			return err
		}
		defer w.Close()

		w.Info("[*] 开始端口扫描，目标: %d，端口: %d，协议: %s", len(targets), len(ports), scanProto)

		// 可选跳过存活探测
		if !scanSkipPing {
			w.Info("[*] 存活探测中...")
			pinger := ping.NewPinger(time.Duration(opts.Timeout)*time.Second, 1, "tcp")
			pingResults := pinger.PingMultiple(targets, opts.Threads)

			var aliveTargets []string
			for _, r := range pingResults {
				if r.Alive {
					aliveTargets = append(aliveTargets, r.Host)
				}
			}
			w.Info("[*] 存活主机: %d/%d", len(aliveTargets), len(targets))
			targets = aliveTargets
		}

		if len(targets) == 0 {
			w.Warn("[!] 没有存活主机")
			return nil
		}

		s := scanner.NewScanner(
			time.Duration(opts.Timeout)*time.Second,
			opts.Threads,
			scanProto,
			opts.Verbose,
		)
		s.Ports = ports

		results := s.ScanMultiple(targets)

		openCount := 0
		for _, r := range results {
			openCount += len(r.Ports)
		}

		w.Info("[*] 扫描完成，发现 %d 个开放端口", openCount)

		return w.WriteHostResults(results)
	},
}

func init() {
	portscanCmd.Flags().StringSliceVarP(&scanTargets, "target", "t", nil, "目标主机 (逗号分隔，支持 CIDR)")
	portscanCmd.Flags().StringVarP(&scanFile, "list", "l", "", "目标列表文件")
	portscanCmd.Flags().StringVarP(&scanPorts, "ports", "p", "top100", "端口范围 (如 80,443,1-1000,top100,top1000,full)")
	portscanCmd.Flags().StringVar(&scanProto, "proto", "tcp", "协议 (tcp/udp)")
	portscanCmd.Flags().BoolVar(&scanSkipPing, "skip-ping", false, "跳过存活探测，直接扫描")

	rootCmd.AddCommand(portscanCmd)
}
