package handler

import (
	"bytes"
	"context"
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
https://marvin.a2tong.com/news/123`)

var apexApi = api.ApexApi{}

func SetUpApex(conf *config.MarvinConfig) {
	apexApi = api.ApexApi{}
	apexApi.Setup(&conf.Apex)
}

func ApexMapQuery() string {
	status, _ := apexApi.GetApexMapStatus(context.Background())
	var buf bytes.Buffer
	err := mapTmp.Execute(&buf, status)
	if err != nil {
		return ""
	}
	return buf.String()
}

func ApexNewsQuery() (string, error) {
	apexNews, err := apexApi.GetApexNews(context.Background())
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = newsTmp.Execute(&buf, apexNews[0])
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
