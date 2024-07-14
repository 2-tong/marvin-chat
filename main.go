package main

import (
	"bytes"
	"context"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"marvin-chat/api"
	"marvin-chat/config"
	"strings"
	"text/template"
	"time"
)

var botApi openapi.OpenAPI = nil
var apexApi *api.ApexApi = nil
var shortApi *api.ShortApi = nil

func main() {
	conf, err := config.LoadConfig("./marvin.yml")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	apexApi = &api.ApexApi{}
	apexApi.Setup(&conf.Apex)
	shortApi = api.NewShortApi(conf.ShortKey)

	botToken := token.BotToken(conf.Marvin.AppID, conf.Marvin.Token)
	botApi = botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second) // 使用NewSandboxOpenAPI创建沙箱环境的实例

	ws, _ := botApi.WS(ctx, nil, "")
	intent := websocket.RegisterHandlers(onGroupMessageIn(), onPrivateMessageIn())
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = botgo.NewSessionManager().Start(ws, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

var mapTemplate = `匹配地图↓↓↓↓↓↓↓↓↓↓
当前地图: {{.BattleRoyale.Current.ChineseName}}
结束时间: {{.BattleRoyale.Current.FixedEndTimeStr}}
下个地图: {{.BattleRoyale.Next.ChineseName}}
结束时间: {{.BattleRoyale.Next.FixedEndTimeStr}}
排位地图↓↓↓↓↓↓↓↓↓↓
当前地图: {{.Ranked.Current.ChineseName}}
结束时间: {{.Ranked.Current.FixedEndTimeStr}}
下个地图: {{.Ranked.Next.ChineseName}}
结束时间: {{.Ranked.Next.FixedEndTimeStr}}`

func apexMapToString(status *api.ApexStatus) string {
	tmpl, err := template.New("test").Parse(mapTemplate)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, status)
	if err != nil {
		return ""
	}
	return buf.String()
}

func replyC2C(origin *dto.Message) {

	replayContent := handleStrContent(origin.Content)

	reMsg := dto.MessageToCreate{
		MsgID:   origin.ID,
		EventID: string(dto.EventC2CMessageCreate),
		MsgType: 0,
		Content: replayContent,
	}
	_, err := botApi.PostC2CMessage(context.Background(), origin.Author.UserOpenid, &reMsg)
	if err != nil {
		return
	}
}

func handleStrContent(msg string) string {
	if strings.Contains(msg, "地图") {
		status, _ := apexApi.GetApexMapStatus(context.Background())
		return apexMapToString(status)
	}
	if strings.Contains(msg, "新闻") || strings.Contains(msg, "咨询") {
		return handleNews()
	}
	return "听不懂思密达😅"
}

func handleNews() string {
	newTmp := `标题：{{.Title}}
内容：{{.ShortDesc}}
{{.Link}}`

	news, err := apexApi.GetApexNews(context.Background())
	if err != nil {
		return "发送内部错误"
	}
	topNews := news[0]
	topNews.Link, err = shortApi.GetLink(context.Background(), topNews.Link)
	if err != nil {
		return "发送内部错误"
	}
	tmpl, err := template.New("newTmp").Parse(newTmp)
	if err != nil {
		return "发送内部错误"
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, topNews)
	if err != nil {
		return ""
	}
	return buf.String()
}

func replyGroup(origin *dto.Message) {
	replayContent := handleStrContent(origin.Content)

	reMsg := dto.MessageToCreate{
		MsgID:   origin.ID,
		EventID: string(dto.EventGroupAtMessageCreate),
		MsgType: 0,
		Content: "\n" + replayContent,
	}
	_, err := botApi.PostGroupMessage(context.Background(), origin.GroupOpenid, &reMsg)
	if err != nil {
		return
	}
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
