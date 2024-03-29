package httpgw

import (
	"errors"
	"github.com/obase/center/httpx"
	"github.com/obase/log"
	"net/http"
)

var plugins map[string]http.HandlerFunc = make(map[string]http.HandlerFunc)

// Note: this method is not thread-safe
func Plugin(name string, f http.HandlerFunc) error {
	if _, ok := plugins[name]; ok {
		return errors.New("duplicate plugin: " + name)
	}
	plugins[name] = f
	return nil
}

func Serve() error {
	return ServeWith(LoadConfig())
}

func ServeWith(config *Config) error {

	config = mergeConfig(config)
	if config.HttpPort == 0 {
		return nil
	}

	mux := http.NewServeMux()

	for _, entry := range config.Entries {
		if len(entry.Plugins) > 0 {
			var chain []http.HandlerFunc
			for _, pname := range entry.Plugins {
				if plugin, ok := plugins[pname]; !ok {
					return errors.New("missing plugin: " + pname)
				} else {
					chain = append(chain, plugin)
				}
			}
			if entry.Https {
				// 启用TLS
				mux.HandleFunc(entry.Source, func(writer http.ResponseWriter, request *http.Request) {
					// 1. 执行middle ware插件
					for _, plugin := range chain {
						plugin(writer, request)
						if IsAbort(request) {
							return // 如果中止立即返回
						}
					}
					// 2. 执行转发
					httpx.ProxyHandlerTLS(entry.Service, entry.Target).ServeHTTP(writer, request)
				})
			} else {
				// 不启用TLS
				mux.HandleFunc(entry.Source, func(writer http.ResponseWriter, request *http.Request) {
					// 1. 执行middle ware插件
					for _, plugin := range chain {
						plugin(writer, request)
						if IsAbort(request) {
							return // 如果中止立即返回
						}
					}
					// 2. 执行转发
					httpx.ProxyHandler(entry.Service, entry.Target).ServeHTTP(writer, request)
				})
			}
		} else {
			if entry.Https {
				// 启用TLS
				mux.Handle(entry.Source, httpx.ProxyHandlerTLS(entry.Service, entry.Target))
			} else {
				// 不启用TLS
				mux.Handle(entry.Source, httpx.ProxyHandler(entry.Service, entry.Target))
			}
		}
	}

	if config.Name != "" {
		registerServiceProxy(mux, config)
	}
	server := &http.Server{
		Handler: mux,
	}

	listner, err := graceListenHttp(config.HttpHost, config.HttpPort, config.HttpKeepAlive)
	if err != nil {
		return err
	}
	defer listner.Close()

	if config.HttpCertFile != "" {
		go func() {
			if err = server.ServeTLS(listner, config.HttpCertFile, config.HttpKeyFile); err != nil {
				log.Error(nil, "server.ServeTLS error: %v", err)
			}
		}()
	} else {
		go func() {
			if err = server.Serve(listner); err != nil {
				log.Error(nil, "server.Serve error: %v", err)
			}
		}()
	}

	// 等待中止信号
	graceShutdownOrRestart(server, listner)
	return nil
}
