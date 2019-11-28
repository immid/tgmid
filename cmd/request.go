package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/immid/tgmid/pkg/base"
	"github.com/immid/tgmid/pkg/tgRpc"
	"github.com/spf13/cobra"
	"net/rpc/jsonrpc"
)

var (
	requestCmd = &cobra.Command{
		Use:   "request [flags] <request_body>",
		Short: "Send a request and wait for a response",
		Long:  "Send a request and wait for a response like http",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			talkTo, err := cmd.Flags().GetInt32("talk-to-user-id")
			if err != nil {
				base.Log("rpc: request:", err)
				fmt.Println(err)
			}
			chatId, err := cmd.Flags().GetInt64("chat-id")
			if err != nil {
				base.Log("rpc: request:", err)
				fmt.Println(err)
			}
			if chatId == 0 {
				chatId = int64(talkTo)
			}
			rpcClient := tgRpc.NewClient()
			defer rpcClient.Close()
			jsonClient := jsonrpc.NewClient(rpcClient)
			var response string
			request := tgRpc.Request{
				ChatId:       chatId,
				TalkToUserId: talkTo,
				Content:      args[0],
			}
			requestJson, err := json.Marshal(request)
			if err != nil {
				base.Log("rpc: request json:", err)
				fmt.Println(err)
			}
			err = jsonClient.Call("ServerHandler.Request", string(requestJson), &response)
			if err != nil {
				base.Log("rpc: request:", err)
				fmt.Println(err)
			}
			fmt.Println(response)
		},
	}
)

func init() {
	requestCmd.Flags().Int32P("talk-to-user-id", "u", 0, "Talk to user id (required)")
	_ = requestCmd.MarkFlagRequired("talk-to-user-id")
	requestCmd.Flags().Int64P("chat-id", "C", 0, `Chat id. if omitted, "talk-to-user-id" will be used`)
	rootCmd.AddCommand(requestCmd)
}
