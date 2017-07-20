package trans_test

import (
	"fmt"
	"reflect"

	"github.com/goodplayer/asa/core/trans"
)

func ExampleTcpLinux() {
	tcpConfig := trans.TcpTransConfig{}
	tcpConfig.UseReuseport = true

	l, err := trans.NewTcpTrans().NewListener("tcp", ":8012", tcpConfig)
	fmt.Println(reflect.TypeOf(l), err)

	for i := 0; i < 3; i++ {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println(conn.RemoteAddr())
		conn.Close()
	}
}

func ExampleMultiTcpListen() {
	tcpConfig := trans.TcpTransConfig{}
	tcpConfig.UseReuseport = true

	_, err := trans.NewTcpTrans().NewListener("tcp", ":8012", tcpConfig)
	if err != nil {
		panic(err)
	}

	_, err = trans.NewTcpTrans().NewListener("tcp", ":8012", tcpConfig)
	if err != nil {
		panic(err)
	}
}
