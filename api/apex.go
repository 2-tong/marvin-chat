package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tencent-connect/botgo/log"
	"marvin-chat/cache"
	"marvin-chat/config"
	"time"
)

const (
	address = "https://api.mozambiquehe.re"

	mapURI      = "/maprotation?version=2"
	mapCacheKey = "mapInfo"
)

type apexMapRequest struct {
	config.ApexConfig
}

type ApexApi struct {
	resCache    *cache.TimeOutCache
	authKey     string
	restyClient *resty.Client // resty client 复用
}

// Setup 生成一个实例
func (o *ApexApi) Setup(apexCfg *config.ApexConfig) *ApexApi {
	o.authKey = apexCfg.AuthKey
	o.resCache = cache.NewCache()
	o.setupClient() // 初始化可复用的 client
	return o
}

// request 每个请求，都需要创建一个 request
func (o *ApexApi) request(ctx context.Context) *resty.Request {
	return o.restyClient.R().SetContext(ctx)
}

// GetApexMapStatus 获取apex地图
func (o *ApexApi) GetApexMapStatus(ctx context.Context) (*ApexStatus, error) {
	get, hit := o.resCache.Get(mapCacheKey)
	if hit {
		return get.(*ApexStatus), nil
	}

	resp, err := o.request(ctx).
		Get(address + mapURI)
	if err != nil {
		return nil, err
	}
	res := ApexStatus{}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	minTime := res.BattleRoyale.Current.End
	rkT := res.Ranked.Current.End

	if minTime > rkT {
		minTime = rkT
	}

	utcTime := time.Unix(minTime, 0)

	o.resCache.Set(mapCacheKey, &res, utcTime)
	return &res, nil
}

// 初始化 client
func (o *ApexApi) setupClient() {
	o.restyClient = resty.New().
		SetLogger(log.DefaultLogger).
		SetDebug(true).
		SetTimeout(3*time.Second).
		SetHeader("Authorization", o.authKey)
	//SetHeader("User-Agent", version.String()).
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

type GameMode struct {
	Current ApexMapInfo `json:"current"`
	Next    ApexMapInfo `json:"next"`
}

type ApexStatus struct {
	BattleRoyale GameMode `json:"battle_royale"`
	Ranked       GameMode `json:"ranked"`
}
