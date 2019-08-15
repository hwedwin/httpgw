package gwproxy

import (
	"net"
	"net/http"
	"os"
	"time"
)

const (
	ABORT_HEADER = "x-abort"
	GRACE_ENV    = "_GRC_"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func Abort(request *http.Request) {
	request.Header.Set(ABORT_HEADER, "true")
}

func IsAbort(request *http.Request) bool {
	return request.Header.Get(ABORT_HEADER) != ""
}

func GetFile(l net.Listener) *os.File {
	file, _ := l.(*net.TCPListener).File()
	return file
}

var PrivateAddress = func(def string) (ret string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return def
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				// 必须私有网段
				if (ip4[0] == 10) || (ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || (ip4[0] == 192 && ip4[1] == 168) {
					return ip4.String()
				}
			}
		}
	}
	return def
}("127.0.0.1")