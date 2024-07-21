package api

import (
	"context"
	"encoding/json"
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

	news = "/news?lang=zh_TW"
)

// Setup 生成一个实例
func (o *ApexApi) Setup(apexCfg *config.MarvinConfig) *ApexApi {
	o.authKey = apexCfg.Apex.AuthKey
	o.resCache = cache.NewCache()
	o.apikey = apexCfg.ShortKey
	o.setupClient() // 初始化可复用的 client
	o.client = resty.New().
		SetLogger(log.DefaultLogger).
		SetDebug(true).
		SetTimeout(3 * time.Second)
	return o
}

// request 每个请求，都需要创建一个 request
func (o *ApexApi) request(ctx context.Context) *resty.Request {
	return o.restyClient.R().SetContext(ctx)
}

// 初始化 client
func (o *ApexApi) setupClient() {
	o.restyClient = resty.New().
		SetLogger(log.DefaultLogger).
		SetTimeout(3*time.Second).
		SetHeader("Authorization", o.authKey)
	//SetHeader("User-Agent", version.String()).
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

// GetApexNews 获取资讯
func (o *ApexApi) GetApexNews(ctx context.Context) ([]ApexNews, error) {
	resp, err := o.request(ctx).
		Get(address + news)
	if err != nil {
		return nil, err
	}
	var res []ApexNews

	err = json.Unmarshal(resp.Body(), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
