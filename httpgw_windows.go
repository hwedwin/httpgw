package httpgw

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func graceListenTCP(addr string) (net.Listener, error) {
	tln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{TCPListener: tln.(*net.TCPListener)}, nil
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
