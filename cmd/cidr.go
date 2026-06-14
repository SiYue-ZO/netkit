package cmd

import (
	"fmt"

	"github.com/SiYue-ZO/netkit/internal/cidrutil"
	"github.com/SiYue-ZO/netkit/internal/output"
	"github.com/spf13/cobra"
)

var cidrCmd = &cobra.Command{
	Use:   "cidr",
	Short: "CIDR 网段工具",
	Long: `CIDR 网段计算工具，支持:
  expand  - 展开 CIDR 为 IP 列表
  info    - 查看 CIDR 网段信息
  aggregate - 聚合多个 CIDR`,
	SilenceUsage: true,
}

var cidrExpandCmd = &cobra.Command{
	Use:   "expand <cidr>",
	Short: "展开 CIDR 为 IP 列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ips, err := cidrutil.ExpandCIDR(args[0])
		if err != nil {
			return err
		}

		w, err := output.NewWriter(opts.Format, opts.Output, opts.NoColor)
		if err != nil {
			return err
		}
		defer w.Close()

		for _, ip := range ips {
			fmt.Fprintln(w, ip)
		}
		return nil
	},
}

var cidrInfoCmd = &cobra.Command{
	Use:   "info <cidr>",
	Short: "查看 CIDR 网段信息",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := cidrutil.ParseCIDR(args[0])
		if err != nil {
			return err
		}

		w, err := output.NewWriter(opts.Format, opts.Output, opts.NoColor)
		if err != nil {
			return err
		}
		defer w.Close()

		return w.WriteCIDRInfo(info)
	},
}

var cidrAggregateCmd = &cobra.Command{
	Use:   "aggregate <cidr1> <cidr2> ...",
	Short: "聚合多个 CIDR",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := cidrutil.AggregateCIDRs(args)
		if err != nil {
			return err
		}

		w, err := output.NewWriter(opts.Format, opts.Output, opts.NoColor)
		if err != nil {
			return err
		}
		defer w.Close()

		for _, cidr := range result {
			fmt.Fprintln(w, cidr)
		}
		return nil
	},
}

func init() {
	cidrCmd.AddCommand(cidrExpandCmd)
	cidrCmd.AddCommand(cidrInfoCmd)
	cidrCmd.AddCommand(cidrAggregateCmd)
	rootCmd.AddCommand(cidrCmd)
}
