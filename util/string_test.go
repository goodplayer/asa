package util_test

import (
	"testing"

	"github.com/goodplayer/asa/util"
)

func TestStripHostAndPort(t *testing.T) {
	t.Log(util.SplitHostAndPort("127.0.0.1:8888"))
	t.Log(util.SplitHostAndPort("[::1]:80"))
	t.Log(util.SplitHostAndPort("127.0.0.1"))
	t.Log(util.SplitHostAndPort("[::1]"))
	t.Log(util.SplitHostAndPort("127.0.0.1:"))
	t.Log(util.SplitHostAndPort("[::1]:"))
}
