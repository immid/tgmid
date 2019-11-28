package cmd

import (
	"fmt"
	"github.com/immid/tgmid/pkg/base"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(base.Version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
