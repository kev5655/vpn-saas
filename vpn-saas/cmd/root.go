package cmd

import (
	"os"

	"vpn-saas/vpn-saas/internal"

	"github.com/spf13/cobra"
)

var verbose bool

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "gptcli",
		Short: "A CLI tool to interact with ChatGPT and manage PDF outputs",
		Long:  `gptcli allows you to ask questions, attach images, and save responses as PDFs. It also provides cleanup functionality for generated PDFs.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			internal.SetVerbose(verbose)
			internal.InitLogger(verbose)
		},
	}
	// Add a persistent flag for verbosity
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Add subcommands
	rootCmd.AddCommand(upCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
