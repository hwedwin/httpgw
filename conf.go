package gwproxy

import (
	"errors"
	"github.com/obase/conf"
	"time"
)

type Entry struct {
	Source  string   // 来源URI(必需且惟一)
	Service string   // 目标服务(必需)
	Target  string   // 目标URI(必需)
	Plugins []string // 服务插件(可选)
	Remark  string   // 备注描述(可选)
	Https   bool     // 是否使用tls
}

type Config struct {
	Name          string        `json:"name" yaml:"name"`                   // 注册服务名,如果没有则不注册
	CheckTimeout  time.Duration `json:"checkTimeout" yaml:"checkTimeout"`   // 注册服务心跳检测超时
	CheckInterval time.Duration `json:"checkInterval" yaml:"checkInterval"` // 注册服务心跳检测间隔
	HttpHost      string        `json:"httpHost" yaml:"httpHost"`           // Http暴露主机,默认首个私有IP
	HttpPort      int           `json:"httpPort" yaml:"httpPort"`           // Http暴露端口, 默认80
	Entries       []*Entry      `json:"entries" yaml:"entries"`             // 代理入口配置
	CertFile      string        `json:"certFile" yaml:"certFile"`           // 启用TLS
	KeyFile       string        `json:"keyFile" yaml:"keyFile"`             // 启用TLS
}

const CKEY = "service"

var ErrConfigNotFound = errors.New("config not found: " + CKEY)

func LoadConfig() (*Config, error) {
	var config *Config
	ok := conf.Scan(CKEY, &config)
	if !ok {
		return nil, ErrConfigNotFound
	}
	return config, nil
}
