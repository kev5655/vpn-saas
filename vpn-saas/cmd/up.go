package cmd

import (
	"github.com/spf13/cobra"
)

var region string
var outputDir string

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Launch VPN server in AWS",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement the logic to launch the VPN server in AWS
		return nil
	},
}

func ApplyUpCmd() *cobra.Command {
	upCmd.Flags().StringVar(&region, "region", "us-east-1", "AWS region to deploy VPN")
	upCmd.Flags().StringVar(&outputDir, "out", "/tmp", "Directory to write WireGuard config")
	return upCmd
}
