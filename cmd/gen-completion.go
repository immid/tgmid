package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	genCompletionCmd = &cobra.Command{
		Use:   "gen-completion",
		Short: "Generate shell completion file",
		Long:  "Generate shell completion file",
		Run: func(cmd *cobra.Command, args []string) {
			shellType, _ := cmd.Flags().GetString("type")
			switch shellType {
			default:
				_, _ = fmt.Fprintln(os.Stderr, "Unsupported shell type", shellType)
			case "bash":
				_ = rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				_ = rootCmd.GenZshCompletion(os.Stdout)
			case "powershell":
				_ = rootCmd.GenPowerShellCompletion(os.Stdout)
			}
		},
	}
)

func init() {
	genCompletionCmd.Flags().StringP("type", "t", "bash", `Shell type, one of "bash", "zsh" or "powershell"`)
	rootCmd.AddCommand(genCompletionCmd)
}
