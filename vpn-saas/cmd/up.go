package cmd

import (
	// "time"
	"vpn-saas/vpn-saas/internal"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Launch VPN server in AWS",
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.StartInstance(internal.Instances["ecs-x300003tygcz-bg-sofia-1"])
		// time.Sleep(30 * time.Second)
		// internal.StopInstance(internal.Instances["ecs-x300003tygcz-bg-sofia-1"])
		return nil
	},
}
