package upstream_test

import (
	"fmt"

	"github.com/goodplayer/asa/core/upstream"
)

func ExampleUpstream_SelectTcp() {
	up := upstream.Upstream{
		List: []string{"tcp://127.0.0.1:8012"},
	}

	conn, err := up.SelectTcp(1)
	fmt.Println(err, conn.TCPConn.RemoteAddr())
}
