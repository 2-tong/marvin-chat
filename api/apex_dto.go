package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"marvin-chat/cache"
	"time"
)

type GameMode struct {
	Current ApexMapInfo `json:"current"`
	Next    ApexMapInfo `json:"next"`
}

type ApexStatus struct {
	BattleRoyale GameMode `json:"battle_royale"`
	Ranked       GameMode `json:"ranked"`
}

type ApexNews struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Img       string `json:"img"`
	ShortDesc string `json:"short_desc"`
}

type ApexApi struct {
	ShortApi
	resCache    *cache.TimeOutCache
	authKey     string
	restyClient *resty.Client // resty client 复用
}

type ApexMapInfo struct {
	Start             int64  `json:"start"`
	End               int64  `json:"end"`
	ReadableDateStart string `json:"readableDate_start"`
	ReadableDateEnd   string `json:"readableDate_end"`
	Map               string `json:"map"`
	Code              string `json:"code"`
	DurationInSecs    int    `json:"DurationInSecs"`
	DurationInMinutes int    `json:"DurationInMinutes"`
	Asset             string `json:"asset"`
	RemainingSecs     int    `json:"remainingSecs,omitempty"`
	RemainingMins     int    `json:"remainingMins,omitempty"`
	RemainingTimer    string `json:"remainingTimer,omitempty"`
}

func (a *ApexMapInfo) FixedEndTimeStr() string {
	return fixTime(a.End)
}

func (a *ApexMapInfo) ChineseName() string {
	switch a.Map {
	case "Broken Moon":
		return "残月"
	case "Kings Canyon":
		return "诸王峡谷"
	case "Olympus":
		return "奥林匹斯"
	case "World's Edge":
		return "世界尽头"
	case "Storm Point":
		return "风暴点"
	default:
		return a.Map
	}
}

func fixTime(t64 int64) string {
	utcTime := time.Unix(t64, 0)

	// 指定目标时区
	targetLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return ""
	}

	// 将UTC时间转换到目标时区
	targetTime := utcTime.In(targetLocation)

	customFormat := "01-02 15:04"
	return targetTime.Format(customFormat)
}
