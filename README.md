# package httpgw
基于HTTP的反向代理组件. 

# Installation
- go get
```
go get -u github.com/obase/conf
go get -u github.com/obase/center
go get -u github.com/obase/log
go get -u github.com/obase/httpgw
```

- go mod
```
go mod edit -require=github.com/obase/httpgw@latest
```

# Configuration
```
# Http的Transport及Client设置
httpx:
   transport:
     dialerTimeout: "30s"
     dialerKeepAlive: "30s"
     maxIdleConns: 10240
     idleConnTimeout: "90s"
     tlsHandshakeTimeout: "10s"
     expectContinueTimeout: "1s"
     maxIdleConnsPerHost: 2048
     responseHeaderTimeout: "5s"
   client:
     timeout: "60s"

# 服务注册中心
center:
  # 代理地址
  address: "10.11.165.44:18500"
  # 缓存超时
  expired: 60

# httpgw反向代理
httpgw:
  name: "testgw"
  # 服务IP,默认本机首个私有地址
  httpHost: "10.11.165.127"
  # 服务端口,默认80
  httpPort: 80
  # 如果启用https
  httpCertFile:
  httpKeyFile:
  # 服务心跳检测超时. 默认5秒
  httpCheckTimeout: "5s"
  # 服务心跳检测间隔. 默认6秒
  httpCheckInterval: "6s"
  # gwproxy代理入口
  entries:
    - {source: "/gw/mul", service: "target", target: "/mul", plugins: ["demo"], https: false, remark: "测试用例"}

```

# Index
- Constants
```
const (
	ABORT_HEADER = "x-abort"
	GRACE_ENV    = "_GRC_"
)
```
- Variables
```
var plugins map[string]http.HandlerFunc = make(map[string]http.HandlerFunc)
```

- func Plugin
```
func Plugin(name string, f http.HandlerFunc) error
```
添加插件

- func Serve
```
func Serve() error 
```
从conf.yml读取配置启动

- func ServeWith
```
func ServeWith(config *Config) error 
```
指定配置启动

- func LoadConfig
```
func LoadConfig() (*Config) 
```
从conf.yml加载配置

# Examples
```

```