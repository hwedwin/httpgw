package gwproxy

import (
	"fmt"
	"github.com/obase/conf"
	"testing"
	"time"
)

type Transport struct {
	DialerTimeout time.Duration `yaml:"dialerTimeout"`
	DialerKeepAlive time.Duration `yaml:"dialerKeepAlive"`
}

func TestListenAndServe(t *testing.T) {
	var d *Transport
	conf.Scan("httpx.transport", &d)
	fmt.Println(d)
}
