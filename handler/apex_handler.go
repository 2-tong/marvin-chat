package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"marvin-chat/api"
	"marvin-chat/config"
	"marvin-chat/image"
	"os"
	"strconv"
	"text/template"

	"github.com/tencent-connect/botgo/log"
)

var mapTmp, _ = template.New("mapTmp").Parse(
	`匹配地图↓↓↓↓↓↓↓↓↓↓
当前地图: {{.BattleRoyale.Current.ChineseName}}
结束时间: {{.BattleRoyale.Current.FixedEndTimeStr}}
下个地图: {{.BattleRoyale.Next.ChineseName}}
结束时间: {{.BattleRoyale.Next.FixedEndTimeStr}}
排位地图↓↓↓↓↓↓↓↓↓↓
当前地图: {{.Ranked.Current.ChineseName}}
结束时间: {{.Ranked.Current.FixedEndTimeStr}}
下个地图: {{.Ranked.Next.ChineseName}}
结束时间: {{.Ranked.Next.FixedEndTimeStr}}`)

var newsTmp, _ = template.New("newsTmp").Parse(
	`标题:{{.Title}}
简讯:{{.ShortDesc}}
{{.Link}}`)

var apexApi = api.ApexApi{}

type ApexImgCard struct {
	MapCode   string `json:"mapCode"`
	MapName   string `json:"mapName"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type ApexImg struct {
	Current ApexImgCard `json:"current"`
	Next    ApexImgCard `json:"next"`
	Rank    ApexImgCard `json:"rank"`
}

func SetUpApex(conf *config.MarvinConfig) {
	apexApi = api.ApexApi{}
	apexApi.Setup(conf)
}

func fileExists(filename string) bool {
	// 使用 os.Stat 获取文件信息
	_, err := os.Stat(filename)

	// 如果 err 为 nil，说明文件存在
	if err == nil {
		return true
	}

	// 如果错误是因为文件不存在，则返回 false
	if os.IsNotExist(err) {
		return false
	}

	// 其他错误，返回 false
	return false
}

func createImage(status *api.ApexStatus, filePath string) error {
	imgData := ApexImg{
		Current: ApexImgCard{
			MapName:   status.BattleRoyale.Current.Map,
			MapCode:   status.BattleRoyale.Current.Code,
			StartTime: status.BattleRoyale.Current.FixedStartTimeStr(),
			EndTime:   status.BattleRoyale.Current.FixedEndTimeStr(),
		},
		Next: ApexImgCard{
			MapName:   status.BattleRoyale.Next.Map,
			MapCode:   status.BattleRoyale.Next.Code,
			StartTime: status.BattleRoyale.Next.FixedStartTimeStr(),
			EndTime:   status.BattleRoyale.Next.FixedEndTimeStr(),
		},
		Rank: ApexImgCard{
			MapName:   status.Ranked.Current.Map,
			MapCode:   status.Ranked.Current.Code,
			StartTime: status.Ranked.Current.FixedStartTimeStr(),
			EndTime:   status.Ranked.Current.FixedEndTimeStr(),
		},
	}

	jsonData, err := json.Marshal(imgData)
	if err != nil {
		log.Error(err)
		return err
	}
	image.DarwImg(config.GetConfig().MapHtml, string(jsonData), filePath)
	return nil
}

func ApexMapQuery(_ string) MsgReply {
	status, err := apexApi.GetApexMapStatus(context.Background())
	if err != nil {
		log.Error(err)
		return MsgReply{Type: "text", Content: ""}
	}
	cfg := config.GetConfig()
	var filePath = cfg.ImgPath + "/" + strconv.FormatInt(status.ExpireT, 10) + ".png"
	if !fileExists(filePath) {
		createImage(status, filePath)
	}

	return MsgReply{
		Type:    "image",
		Content: cfg.ImgUrlPrefix + strconv.FormatInt(status.ExpireT, 10) + ".png",
	}
}

func ApexNewsQuery(_ string) MsgReply {
	apexNews, err := apexApi.GetApexNews(context.Background())
	if err != nil {
		log.Error(err)
		return MsgReply{Type: "text", Content: ""}
	}
	news := apexNews[0]
	news.Link, _ = apexApi.GetLink(context.Background(), news.Link)
	var buf bytes.Buffer
	err = newsTmp.Execute(&buf, news)
	if err != nil {
		log.Error(err)
		return MsgReply{Type: "text", Content: ""}
	}
	return MsgReply{Type: "text", Content: buf.String()}
}

func init() {
	RegisterSimpleMsgHandler("地图", ApexMapQuery)
	RegisterSimpleMsgHandler("新闻", ApexNewsQuery)
	RegisterSimpleMsgHandler("资讯", ApexNewsQuery)
}
