package gwproxy

import (
	"net/http"
	"time"
)

type Entry struct {
	Source  string   // 来源URI(必需且惟一)
	Service string   // 目标服务(必需)
	Target  string   // 目标URI(必需)
	Plugins []string // 服务插件(可选)
	Remark  string   // 备注描述(可选)
}

type Config struct {
	Name          string        // 注册服务名,如果没有则不注册
	CheckTimeout  time.Duration // 注册服务心跳检测超时
	CheckInterval time.Duration // 注册服务心跳检测间隔
	HttpHost      string        // Http暴露主机,默认首个私有IP
	HttpPort      int           // Http暴露端口, 默认80
	Entries       []*Entry      // 代理入口配置
}

func Plugin(name string, f http.HandlerFunc) {

}

func ListenAndServe() error {
	return nil
}
