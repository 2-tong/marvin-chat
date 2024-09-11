module marvin-chat

go 1.22.4

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-resty/resty/v2 v2.6.0
	github.com/tencent-connect/botgo v0.1.6
	gopkg.in/yaml.v3 v3.0.0
)

replace github.com/tencent-connect/botgo => ../botgo

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/tidwall/gjson v1.9.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/net v0.17.0 // indirect
)
