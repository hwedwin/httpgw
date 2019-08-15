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
	ProxyName     string        `json:"proxyName" bson:"proxyName" yaml:"proxyName"`             // 注册服务名,如果没有则不注册
	CheckTimeout  time.Duration `json:"checkTimeout" bson:"checkTimeout" yaml:"checkTimeout"`    // 注册服务心跳检测超时
	CheckInterval time.Duration `json:"checkInterval" bson:"checkInterval" yaml:"checkInterval"` // 注册服务心跳检测间隔
	ProxyHost     string        `json:"proxyHost" bson:"proxyHost" yaml:"proxyHost"`             // Http暴露主机,默认首个私有IP
	ProxyPort     int           `json:"proxyPort" bson:"proxyPort" yaml:"proxyPort"`             // Http暴露端口, 默认80
	ProxyCertFile string        `json:"proxyCertFile" bson:"proxyCertFile" yaml:"proxyCertFile"` // 启用TLS
	ProxyKeyFile  string        `json:"proxyKeyFile" bson:"proxyKeyFile" yaml:"proxyKeyFile"`    // 启用TLS
	Entries       []*Entry      `json:"entries" json:"entries" yaml:"entries"`                   // 代理入口配置
}

const CKEY = "service"

var ErrConfigNotFound = errors.New("missing map config: service")

func LoadConfig() (*Config) {
	var config *Config
	ok := conf.Scan(CKEY, &config)
	if !ok {
		return nil
	}
	return config
}
