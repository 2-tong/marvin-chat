package main

import (
	"context"
	"log"
	"marvin-chat/config"
	"marvin-chat/handler"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

var botApi openapi.OpenAPI = nil

func main() {
	conf := config.GetConfig()
	ctx := context.Background()
	handler.SetUpApex(conf)
	botToken := token.NewQQBotTokenSource(&token.QQBotCredentials{AppID: conf.Marvin.AppID, AppSecret: conf.Marvin.AppSecret})
	botApi = botgo.NewSandboxOpenAPI(conf.Marvin.AppID, botToken).WithTimeout(3 * time.Second) // 使用NewSandboxOpenAPI创建沙箱环境的实例

	ws, _ := botApi.WS(ctx, nil, "")
	intent := websocket.RegisterHandlers(onGroupMessageIn(), onPrivateMessageIn())
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	var err error = nil
	if err = botgo.NewSessionManager().Start(ws, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// 统一处理消息发送的函数
func handleReply(origin *dto.Message, isGroup bool) {
	handler.HandleTextMsg(origin.Content, func(reply handler.MsgReply) {
		reMsg := dto.MessageToCreate{
			MsgID: origin.ID,
		}

		if isGroup {
			reMsg.EventID = string(dto.EventGroupAtMessageCreate)
		} else {
			reMsg.EventID = string(dto.EventC2CMessageCreate)
		}

		switch reply.Type {
		case "image":
			richMsg := dto.RichMediaMessage{
				FileType:   1,
				URL:        reply.Content,
				SrvSendMsg: false,
			}

			var fileBack *dto.Message
			var errImg error
			if isGroup {
				fileBack, errImg = botApi.PostGroupMessage(context.Background(), origin.GroupID, &richMsg)
			} else {
				fileBack, errImg = botApi.PostC2CMessage(context.Background(), origin.Author.ID, &richMsg)
			}
			if errImg != nil {
				return
			}
			reMsg.MsgType = dto.RichMediaMsg
			reMsg.Media = &dto.MediaInfo{FileInfo: fileBack.FileInfo}
		case "text":
			reMsg.MsgType = dto.TextMsg
			reMsg.Content = reply.Content
		}

		var err error
		if isGroup {
			_, err = botApi.PostGroupMessage(context.Background(), origin.GroupID, &reMsg)
		} else {
			_, err = botApi.PostC2CMessage(context.Background(), origin.Author.ID, reMsg)
		}
		if err != nil {
			return
		}
	})
}

func replyC2C(origin *dto.Message) {
	handleReply(origin, false)
}

func replyGroup(origin *dto.Message) {
	handleReply(origin, true)
}

func onPrivateMessageIn() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		replyC2C((*dto.Message)(data))
		return nil
	}
}

func onGroupMessageIn() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		replyGroup((*dto.Message)(data))
		return nil
	}
}
