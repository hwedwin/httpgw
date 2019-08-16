package httpgw

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func graceListenHttp(host string, port int, keepAlivePeriod time.Duration) (net.Listener, error) {
	tln, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{TCPListener: tln.(*net.TCPListener), KeepAlivePeriod: keepAlivePeriod}, nil
}

func graceShutdownOrRestart(httpServer *http.Server, httpListener net.Listener) {
	sch := make(chan os.Signal, 1)
	defer signal.Stop(sch)

	signal.Notify(sch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-sch
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			httpServer.Shutdown(context.Background())
			return
		}
	}
}
