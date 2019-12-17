package tgRpc

import (
	"github.com/Arman92/go-tdlib"
	"github.com/immid/tgmid/pkg/base"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

type MsgMonitor struct {
	ChatId       int64
	SenderId     int32
	ReplyToMsgId int64
	OldMsgId     int64
	CallbackData chan *tdlib.Message
}

func (mon MsgMonitor) Match(message *tdlib.Message, oldMsgId int64) bool {
	base.LogVerbose("rpc: monitor: matching", message.ID)
	if mon.ChatId != message.ChatID {
		base.LogVerbose("rpc: monitor: chat id: expect:", mon.ChatId, "got:", message.ChatID)
		return false
	}
	if mon.SenderId != message.SenderUserID {
		base.LogVerbose("rpc: monitor: sender id: expect:", mon.SenderId, "got:", message.SenderUserID)
		return false
	}
	if mon.ReplyToMsgId != message.ReplyToMessageID {
		base.LogVerbose("rpc: monitor: reply to id: expect:", mon.ReplyToMsgId, "got:", message.ReplyToMessageID)
		return false
	}
	if mon.OldMsgId != oldMsgId {
		base.LogVerbose("rpc: monitor: old msg id: expect:", mon.OldMsgId, "got:", oldMsgId)
		return false
	}
	base.LogVerbose("rpc: monitor: matched", message.ID)
	return true
}

func (mon MsgMonitor) Process(message *tdlib.Message) {
	mon.CallbackData <- message
}

type ServerHandler struct {
	client   *tdlib.Client
	monitors map[int64]*MsgMonitor
	counter  *base.Counter
}

func AddMonitor(handler *ServerHandler, mon *MsgMonitor) int64 {
	index := handler.counter.Next()
	handler.monitors[index] = mon
	base.LogVerbose("rpc: monitor added:", index)
	return index
}

func RemoveMonitor(handler *ServerHandler, index int64) {
	delete(handler.monitors, index)
	base.LogVerbose("rpc: monitor removed:", index)
}

func ListenMessages(handler *ServerHandler, messageType tdlib.TdMessage) {
	receiver := handler.client.AddEventReceiver(messageType, func(msg *tdlib.TdMessage) bool {
		return true
	}, 100)
	base.Log("rpc: msg listener registered:", messageType.MessageType(), "\n")
	for tdMessage := range receiver.Chan {
		var message *tdlib.Message
		var oldMsgId int64 = 0
		switch tdMessage.(type) {
		default:
			continue
		case *tdlib.UpdateNewMessage:
			message = tdMessage.(*tdlib.UpdateNewMessage).Message
		case *tdlib.UpdateMessageSendSucceeded:
			sentMsg := tdMessage.(*tdlib.UpdateMessageSendSucceeded)
			message = sentMsg.Message
			oldMsgId = sentMsg.OldMessageID
		}
		base.LogVerbose("rpc: msg listener:", messageType.MessageType(), message.ID, "from:", message.SenderUserID, "\n")
		for monIndex, msgMon := range handler.monitors {
			base.LogVerbose("rpc: msg type", messageType.MessageType(), "monitor found:", msgMon)
			if msgMon.Match(message, oldMsgId) {
				base.LogVerbose("rpc: msg type", messageType.MessageType(), "monitor processing:", message.ID)
				msgMon.Process(message)
				RemoveMonitor(handler, monIndex)
			}
		}
	}
}

func NewClient() net.Conn {
	socketType := base.GetConfig().GetString("rpc.socket_type")
	address := base.GetConfig().GetString("rpc.address")
	timeout := base.GetConfig().GetDuration("rpc.timeout")
	client, err := net.DialTimeout(socketType, address, timeout*time.Second)
	if err != nil {
		base.Log("rpc: dial: ", err)
	}
	return client
}

func Serve(client *tdlib.Client) {
	socketType := base.GetConfig().GetString("rpc.socket_type")
	address := base.GetConfig().GetString("rpc.address")
	server := rpc.NewServer()
	if socketType == "unix" {
		_ = os.Remove(address)
	}
	rpcSocket, err := net.Listen(socketType, address)
	if err != nil {
		base.Log("rpc: listen:", err)
		return
	}
	defer rpcSocket.Close()

	base.Log("rpc: listening on", socketType, address)
	handler := &ServerHandler{
		client:   client,
		monitors: make(map[int64]*MsgMonitor),
		counter:  base.NewCounter(0),
	}
	err = server.Register(handler)
	if err != nil {
		base.Log("rpc: register:", err)
		return
	}
	go ListenMessages(handler, &tdlib.UpdateMessageSendSucceeded{})
	go ListenMessages(handler, &tdlib.UpdateNewMessage{})
	for {
		conn, err := rpcSocket.Accept()
		if err != nil {
			base.Log("rpc: accept:", err)
			return
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
