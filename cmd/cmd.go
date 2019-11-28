package cmd

import (
	"fmt"
	"github.com/immid/tgmid/pkg/base"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tgmid",
		Short: "Telegram Middleman",
		Long:  `Telegram client as middleman`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func Execute() {
	rootCmd.PersistentFlags().CountVarP(&base.Verbosity, "verbose", "v", "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&base.ConfigDir, "config-dir", "c", "./configs/", "Set config dir")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
