package api

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/tencent-connect/botgo/log"
	"strings"
	"time"
)

// Response represents the overall response structure
type Response struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

// Data represents the data part of the response
type Data struct {
	Group       Group  `json:"group"`
	Link        Link   `json:"link"`
	ReportURL   string `json:"report_url"`
	NLinksToday int    `json:"n_links_today"`
}

// Group represents the group information
type Group struct {
	Name string `json:"name"`
	Sid  string `json:"sid"`
}

// Link represents the link information
type Link struct {
	Name      string `json:"name"`
	OriginURL string `json:"origin_url"`
	URL       string `json:"url"`
}

type Req struct {
	Apikey       string `json:"apikey"`
	Domain       string `json:"domain"`
	OriginUrl    string `json:"origin_url"`
	GroupSid     string `json:"group_sid"`
	Report       bool   `json:"report"`
	Webhook      bool   `json:"webhook"`
	WebhookScene string `json:"webhook_scene"`
}

type ShortApi struct {
	apikey string
	client *resty.Client
}

func NewShortApi(apikey string) *ShortApi {
	api := &ShortApi{
		apikey: apikey,
		client: resty.New().
			SetLogger(log.DefaultLogger).
			SetDebug(true).
			SetTimeout(3 * time.Second),
	}

	return api
}

func (o *ShortApi) GetLink(ctx context.Context, originLink string) (string, error) {
	res, err := o.client.R().SetContext(ctx).
		SetResult(&Response{}).
		SetBody(&Req{
			Apikey:    o.apikey,
			OriginUrl: originLink,
		}).
		Post("https://api.xiaomark.com/v1/link/create")
	if err != nil {
		return "", err
	}
	url := res.Result().(*Response).Data.Link.URL
	parts := strings.Split(url, "/")
	if len(parts) > 3 {
		shortCode := parts[3]
		url = "https://marvin.a2tong.com/news/" + shortCode
	} else {
		return "", err
	}
	return url, nil
}
