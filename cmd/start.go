package cmd

import (
	"fmt"
	"github.com/immid/tgmid/pkg/base"
	"github.com/immid/tgmid/pkg/telegram"
	"github.com/immid/tgmid/pkg/tgRpc"
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the telegram client",
		Long:  "Start the telegram client and listening to instructions",
		Run: func(cmd *cobra.Command, args []string) {
			client, ready := telegram.Login()
			processor := base.GetConfig().GetString("processor.cmd")
			if !ready {
				base.Log("Not ready")
				return
			}
			go tgRpc.Serve(client)
			rawUpdates := client.GetRawUpdatesChannel(100)
			for update := range rawUpdates {
				jsonData := string(update.Raw)
				fmt.Println(base.LocalExec(processor, jsonData), "\n")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
}

