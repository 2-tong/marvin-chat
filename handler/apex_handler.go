package handler

import (
	"bytes"
	"context"
	"github.com/tencent-connect/botgo/log"
	"marvin-chat/api"
	"marvin-chat/config"
	"text/template"
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

func SetUpApex(conf *config.MarvinConfig) {
	apexApi = api.ApexApi{}
	apexApi.Setup(conf)
}

func ApexMapQuery(_ string) string {
	status, _ := apexApi.GetApexMapStatus(context.Background())
	var buf bytes.Buffer
	err := mapTmp.Execute(&buf, status)
	if err != nil {
		log.Error(err)
		return ""
	}
	return buf.String()
}

func ApexNewsQuery(_ string) string {
	apexNews, err := apexApi.GetApexNews(context.Background())
	if err != nil {
		log.Error(err)
		return ""
	}
	news := apexNews[0]
	news.Link, _ = apexApi.GetLink(context.Background(), news.Link)
	var buf bytes.Buffer
	err = newsTmp.Execute(&buf, news)
	if err != nil {
		log.Error(err)
		return ""
	}
	return buf.String()
}

func init() {
	RegisterSimpleMsgHandler("地图", ApexMapQuery)
	RegisterSimpleMsgHandler("新闻", ApexNewsQuery)
	RegisterSimpleMsgHandler("资讯", ApexNewsQuery)
}
