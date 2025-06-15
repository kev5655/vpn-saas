package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var region string
var outputDir string

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Launch VPN server in AWS",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Generate keypair
		privKeyCmd := exec.Command("wg", "genkey")
		privKeyOut, err := privKeyCmd.Output()
		if err != nil {
			return err
		}
		pubKeyCmd := exec.Command("wg", "pubkey")
		pubKeyCmd.Stdin = bytes.NewReader(privKeyOut)
		pubKeyOut, err := pubKeyCmd.Output()
		if err != nil {
			return err
		}

		// 2. Set env for CDK and deploy
		os.Setenv("CLIENT_PUBKEY", string(pubKeyOut))
		os.Setenv("CDK_REGION", region)
		cdk := exec.Command("cdk", "deploy", "VpnStack", "--require-approval", "never")
		cdk.Stdout = os.Stdout
		cdk.Stderr = os.Stderr
		if err := cdk.Run(); err != nil {
			return err
		}

		// 3. TODO: describe ECS service, fetch public IP (via AWS SDK / AWS CLI)
		fmt.Println("✅ VPN deployed. Config at:", filepath.Join(outputDir, region+"-wg0.conf"))
		return nil
	},
}

func init() {
	upCmd.Flags().StringVar(&region, "region", "us-east-1", "AWS region to deploy VPN")
	upCmd.Flags().StringVar(&outputDir, "out", "/tmp", "Directory to write WireGuard config")
	rootCmd.AddCommand(upCmd)
}
