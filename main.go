package main

import (
	"context"
	"fmt"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"marvin-chat/config"
	"time"
)

var api openapi.OpenAPI = nil

func main() {
	conf, err := config.LoadConfig("./marvin.yml")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	botToken := token.BotToken(conf.AppID, conf.Token)
	api = botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second) // 使用NewSandboxOpenAPI创建沙箱环境的实例

	ws, _ := api.WS(ctx, nil, "")
	intent := websocket.RegisterHandlers(onGroupMessageIn(), onPrivateMessageIn())
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = botgo.NewSessionManager().Start(ws, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

func replyC2C(origin *dto.Message) {
	reMsg := dto.MessageToCreate{
		MsgID:   origin.ID,
		EventID: string(dto.EventC2CMessageCreate),
		MsgType: 0,
		Content: "回复：" + origin.Content,
	}
	message, err := api.PostC2CMessage(context.Background(), origin.Author.UserOpenid, &reMsg)
	if err != nil {
		return
	}
	fmt.Println(message)
}

func replyGroup(origin *dto.Message) {
	reMsg := dto.MessageToCreate{
		MsgID:   origin.ID,
		EventID: string(dto.EventGroupAtMessageCreate),
		MsgType: 0,
		Content: "回复：" + origin.Content,
	}
	message, err := api.PostGroupMessage(context.Background(), origin.GroupOpenid, &reMsg)
	if err != nil {
		return
	}
	fmt.Println(message)
}

func onPrivateMessageIn() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		replyC2C((*dto.Message)(data))
		return nil
	}
}

func onGroupMessageIn() event.GroupAtMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupAtMessageData) error {
		replyGroup((*dto.Message)(data))
		return nil
	}
}
