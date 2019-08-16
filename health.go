package gwproxy

import (
	"fmt"
	"github.com/obase/center"
	"github.com/obase/log"
	"net/http"
	"strconv"
)

var OK = []byte("OK")

func CheckHttpHealth(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write(OK)
}

func registerServiceProxy(mux *http.ServeMux, conf *Config) {
	defer log.Flush()
	mux.HandleFunc("/health", CheckHttpHealth)

	suffix := "@" + conf.ProxyHost + ":" + strconv.Itoa(conf.ProxyPort)
	myname := center.ProxyName(conf.Name)
	regs := &center.Service{
		Id:   myname + suffix,
		Kind: "proxy",
		Name: myname,
		Host: conf.ProxyHost,
		Port: conf.ProxyPort,
	}

	chks := &center.Check{
		Type:     "http",
		Target:   fmt.Sprintf("http://%s:%v/health", conf.ProxyHost, conf.ProxyPort),
		Timeout:  conf.ProxyCheckTimeout,
		Interval: conf.ProxyCheckInterval,
	}

	if err := center.Register(regs, chks); err == nil {
		log.Info(nil, "register service success, %v", *regs)
	} else {
		log.Error(nil, "register service error, %v, %v", *regs, err)
	}
}
