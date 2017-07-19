package trans_test

import (
	"fmt"
	"reflect"
	"time"

	"github.com/goodplayer/asa/core/trans"
)

func ExampleTcpLinux() {
	tcpConfig := trans.TcpTransConfig{}
	tcpConfig.UseReuseport = true

	l, _ := trans.NewTcpTrans().NewListener("tcp", ":8012", tcpConfig)
	fmt.Println(reflect.TypeOf(l))

	time.Sleep(1 * time.Hour)
}
