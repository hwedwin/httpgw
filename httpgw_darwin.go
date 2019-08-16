package httpgw

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var flag = os.Getenv(GRACE_ENV)

func graceListenTCP(addr string) (net.Listener, error) {
	if flag != "" {
		file := os.NewFile(3, "")
		defer file.Close()
		if grpcListner, err = net.FileListener(file); err != nil {
			log.Error(nil, "FileListener error: %v", err)
		}
		return grpcListner, err
	}
	tln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return tcpKeepAliveListener{TCPListener: tln.(*net.TCPListener)}, nil
}

func graceShutdownOrRestart(httpServer *http.Server, httpListener net.Listener) {
	sch := make(chan os.Signal, 1)
	defer signal.Stop(sch)

	signal.Notify(sch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-sch
		switch sig {
		case syscall.SIGUSR2:
			var args []string
			// 设置重启标志及参数
			if len(os.Args) > 1 {
				args = os.Args[1:]
			}
			// 执行重启命令
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), GRACE_ENV+"=1") // 拼加标志
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.ExtraFiles = []*os.File{GetFile(httpListener)}
			if err := cmd.Start(); err != nil {
				log.Error(nil, "restart error: %v", err)
			}
			fallthrough
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			httpServer.Shutdown(context.Background())
			return
		}
	}
}
