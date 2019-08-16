package httpgw

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

	shost := conf.HttpHost
	sport := strconv.Itoa(conf.HttpPort)
	if shost == "" {
		shost = PrivateAddress
	}

	suffix := "@" + shost + ":" + sport
	myname := center.ProxyName(conf.Name)
	regs := &center.Service{
		Id:   myname + suffix,
		Kind: "gw",
		Name: myname,
		Host: shost,
		Port: conf.HttpPort,
	}
	chks := &center.Check{
		Type:     "http",
		Target:   fmt.Sprintf("http://%s:%v/health", shost, sport),
		Timeout:  conf.HttpCheckTimeout,
		Interval: conf.HttpCheckInterval,
	}

	if err := center.Register(regs, chks); err == nil {
		log.Info(nil, "register httpgw success, %v", *regs)
	} else {
		log.Error(nil, "register httpgw error, %v, %v", *regs, err)
	}
}
