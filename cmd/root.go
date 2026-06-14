package cmd

import (
	"fmt"
	"os"

	"github.com/SiYue-ZO/netkit/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	opts    types.ScanOptions
)

var rootCmd = &cobra.Command{
	Use:   "netkit",
	Short: "NetKit - 网络安全工具集",
	Long: `NetKit 是一个用 Go 编写的多功能网络侦察与安全检查工具。

面向网络管理员日常运维、安全工程师资产梳理、开发者学习网络编程。

子命令:
  ping       主机存活探测 (ICMP/TCP)
  portscan   端口扫描 (CONNECT/SYN/UDP)
  cidr       CIDR 网段工具
  dns        DNS 枚举与查询
  fingerprint 服务指纹识别
  web        Web 侦察
  ssl        SSL/TLS 安全检查
  trace      路由追踪
  whois      Whois 查询

示例:
  netkit ping -t 192.168.1.0/24
  netkit portscan -t 192.168.1.1 -p 1-65535
  netkit cidr expand 192.168.1.0/24`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "配置文件路径")
	rootCmd.PersistentFlags().IntVarP(&opts.Timeout, "timeout", "T", 5, "超时时间(秒)")
	rootCmd.PersistentFlags().IntVarP(&opts.Threads, "threads", "c", 50, "并发线程数")
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "详细输出")
	rootCmd.PersistentFlags().StringVarP(&opts.Output, "output", "o", "", "输出文件路径")
	rootCmd.PersistentFlags().StringVarP(&opts.Format, "format", "f", "table", "输出格式 (table/json/csv)")
	rootCmd.PersistentFlags().BoolVar(&opts.NoColor, "no-color", false, "禁用彩色输出")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("netkit")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.netkit")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && opts.Verbose {
		fmt.Fprintln(os.Stderr, "使用配置文件:", viper.ConfigFileUsed())
	}
}
