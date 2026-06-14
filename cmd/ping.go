package cmd

import (
	"fmt"
	"time"

	"github.com/SiYue-ZO/netkit/internal/output"
	"github.com/SiYue-ZO/netkit/internal/ping"
	"github.com/spf13/cobra"
)

var (
	pingMethod  string
	pingCount   int
	pingTargets []string
	pingFile    string
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "主机存活探测 (ICMP/TCP)",
	Long: `主机存活探测工具，支持:
  - ICMP Ping (需要管理员/root 权限)
  - TCP Ping (尝试连接常见端口)
  - 自动模式 (先 ICMP，失败后 TCP)
  - 支持 CIDR/文件输入，并发探测

示例:
  netkit ping -t 192.168.1.1
  netkit ping -t 192.168.1.0/24 -c 100
  netkit ping -t 192.168.1.1,192.168.1.2 --method tcp
  netkit ping -l hosts.txt --method icmp`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 解析目标
		targets, err := ping.ResolveTargets(pingTargets, pingFile)
		if err != nil {
			return err
		}
		if len(targets) == 0 {
			return fmt.Errorf("未指定目标，使用 -t 或 -l 指定")
		}

		w, err := output.NewWriter(opts.Format, opts.Output, opts.NoColor)
		if err != nil {
			return err
		}
		defer w.Close()

		w.Info("[*] 开始存活探测，目标数量: %d，方法: %s", len(targets), pingMethod)

		pinger := ping.NewPinger(
			time.Duration(opts.Timeout)*time.Second,
			pingCount,
			pingMethod,
		)

		results := pinger.PingMultiple(targets, opts.Threads)

		alive := 0
		for _, r := range results {
			if r.Alive {
				alive++
			}
		}

		w.Info("[*] 探测完成: %d/%d 存活", alive, len(targets))

		return w.WriteHostResults(results)
	},
}

func init() {
	pingCmd.Flags().StringSliceVarP(&pingTargets, "target", "t", nil, "目标主机 (逗号分隔，支持 CIDR)")
	pingCmd.Flags().StringVarP(&pingFile, "list", "l", "", "目标列表文件")
	pingCmd.Flags().StringVar(&pingMethod, "method", "auto", "探测方法 (icmp/tcp/auto)")
	pingCmd.Flags().IntVar(&pingCount, "count", 1, "每个目标探测次数")

	rootCmd.AddCommand(pingCmd)
}
