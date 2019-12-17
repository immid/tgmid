package tgRpc

import (
	"encoding/json"
	"github.com/Arman92/go-tdlib"
	"github.com/immid/tgmid/pkg/base"
)

type Request struct {
	ChatId       int64  `json:"chat_id"`
	TalkToUserId int32  `json:"talk_to_user_id"`
	ReplyToMsgId int64  `json:"reply_to_msg_id"`
	Content      string `json:"content"`
	NoResponse   bool   `json:"no_response"`
}

func (handler *ServerHandler) Request(requestJson string, responseJson *string) error {
	client := handler.client
	var request Request
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		base.Log("rpc: handler: request: ", err)
		*responseJson = `{"error": "bad request json"}`
		return nil
	}
	base.LogVerbose(request)
	msg := tdlib.NewInputMessageText(tdlib.NewFormattedText(request.Content, nil), true, true)
	pendingMsg, err := client.SendMessage(request.ChatId, request.ReplyToMsgId, false, true, nil, msg)
	if pendingMsg != nil {
		base.LogVerbose("rpc: handler: pending msg id:", pendingMsg.ID)
	} else {
		return err
	}
	sentMonitor := MsgMonitor{
		ChatId:       request.ChatId,
		SenderId:     pendingMsg.SenderUserID,
		ReplyToMsgId: pendingMsg.ReplyToMessageID,
		OldMsgId:     pendingMsg.ID,
		CallbackData: make(chan *tdlib.Message),
	}
	AddMonitor(handler, &sentMonitor)
	base.LogVerbose("rpc: handler: confirming msg id:", pendingMsg.ID)
	sentMsg := <-sentMonitor.CallbackData
	base.LogVerbose("rpc: handler: sent msg:", sentMsg.ID)
	if request.NoResponse {
		return nil
	}
	responseMonitor := MsgMonitor{
		ChatId:       request.ChatId,
		SenderId:     request.TalkToUserId,
		ReplyToMsgId: sentMsg.ID,
		CallbackData: make(chan *tdlib.Message),
	}
	AddMonitor(handler, &responseMonitor)
	base.LogVerbose("rpc: handler: waiting response for msg:", sentMsg.ID)
	responseMsg := <-responseMonitor.CallbackData
	base.LogVerbose("rpc: handler: msg:", sentMsg.ID, "got response:", responseMsg.ID)
	*responseJson = responseMsg.Content.(*tdlib.MessageText).Text.Text
	base.LogVerbose2("rpc: response:", *responseJson)
	return err
}
