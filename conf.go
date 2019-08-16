package httpgw

import (
	"errors"
	"github.com/obase/conf"
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
	Name              string   `json:"name" bson:"name" yaml:"name"`                                           // 注册服务名,如果没有则不注册
	HttpCheckTimeout  string   `json:"proxyCheckTimeout" bson:"proxyCheckTimeout" yaml:"proxyCheckTimeout"`    // 注册服务心跳检测超时
	HttpCheckInterval string   `json:"proxyCheckInterval" bson:"proxyCheckInterval" yaml:"proxyCheckInterval"` // 注册服务心跳检测间隔
	HttpHost          string   `json:"proxyHost" bson:"proxyHost" yaml:"proxyHost"`                            // Http暴露主机,默认首个私有IP
	HttpPort          int      `json:"proxyPort" bson:"proxyPort" yaml:"proxyPort"`                            // Http暴露端口, 默认80
	HttpCertFile      string   `json:"proxyCertFile" bson:"proxyCertFile" yaml:"proxyCertFile"`                // 启用TLS
	HttpKeyFile       string   `json:"proxyKeyFile" bson:"proxyKeyFile" yaml:"proxyKeyFile"`                   // 启用TLS
	Entries           []*Entry `json:"entries" json:"entries" yaml:"entries"`                                  // 代理入口配置
}

const (
	CKEY = "httpgw"
)

var ErrConfigNotFound = errors.New("missing map config: httpgw")

func LoadConfig() (*Config) {
	var config *Config
	ok := conf.Scan(CKEY, &config)
	if !ok {
		return nil
	}
	return config
}

func mergeConfig(conf *Config) *Config {
	if conf == nil {
		conf = &Config{}
	}

	// 补充默认逻辑
	if conf.HttpHost == "" {
		conf.HttpHost = PrivateAddress
	}
	if conf.HttpCheckTimeout == "" {
		conf.HttpCheckTimeout = "5s"
	}
	if conf.HttpCheckInterval == "" {
		conf.HttpCheckInterval = "6s"
	}
	return conf
}
