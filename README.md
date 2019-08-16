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

# gwproxy代理设置
service:
  name: "gwproxy"
  # 服务IP,默认本机首个私有地址
  proxyHost: "10.11.165.44"
  # 服务端口,默认80
  proxyPort:
  # 如果启用https
  proxyCertFile:
  proxyKeyFile:
  # 服务心跳检测超时. 默认5秒
  proxyCheckTimeout: "5s"
  # 服务心跳检测间隔. 默认6秒
  proxyCheckInterval: "6s"
  # consul代理地址. 默认127.0.0.1:8500
  # 形式1, 指定代理地址
  # center: "127.0.0.1:8500"
  # 形式2, 指定代理地址与缓存超时
  center:
    # 代理地址
    address: "10.11.165.44:18500"
    # 缓存超时
    expired: 60
  # 形式3, 指定静态代理,一般用于本地测试
  # center:
  #   configs: {"abc":["1.2.3.4:8000"]}

# gwproxy代理入口
entries:
  -
    # 来源URI(必需), 所有入口来源必须惟一
    source: "/gamegw/jx3/web/get-corps-top200"
    # 目标服务
    service(必需): "jx3mine"
    # 目标URI
    target(必需): "/jx3mine/get-corps-top200"
    # 中间件服务(可选), 具体参照common的midware
    plugins: ["VerifySign","VerifyTs"]
    # 描述备注(可选)
    remark: "这是一个测试案例"
    # 是否使用https
    https: false
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