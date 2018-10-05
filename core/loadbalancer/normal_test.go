package loadbalancer

import (
	"testing"

	"github.com/goodplayer/asa/api"
)

func TestLoadBalancer_Add(t *testing.T) {
	lb := NewLoadBalancer(api.RANDOM)
	lbImpl := lb.(*LoadBalancer)
	lbImpl.Add(1, new(api.Endpoint))
	if _, ok := lbImpl.mapping[1]; !ok {
		t.Fatal("mapping does not contains key 1")
	}
	if lbImpl.cnt != 1 {
		t.Fatal("cnt is not 1")
	}
	if lbImpl.data[0][0] == nil {
		t.Fatal("key 1 data is nil")
	}
}

func TestLoadBalancer_Add_size256(t *testing.T) {
	lb := NewLoadBalancer(api.RANDOM)
	lbImpl := lb.(*LoadBalancer)
	for i := 0; i < 256+1; i++ {
		lbImpl.Add(i, new(api.Endpoint))
	}
	if _, ok := lbImpl.mapping[256+1-1]; !ok {
		t.Fatal("mapping does not contains key 256+1")
	}
	if lbImpl.cnt != 256+1 {
		t.Fatal("cnt is not 256+1")
	}
	if lbImpl.data[1][0] == nil {
		t.Fatal("key 256+1 data is nil")
	}
	for i, v := range lbImpl.data[0] {
		if v == nil {
			t.Fatal("index ", i, "is nil")
		}
	}
}

func TestLoadBalancer_Pick(t *testing.T) {
	lb := NewLoadBalancer(api.RANDOM)
	lbImpl := lb.(*LoadBalancer)
	for i := 0; i < 256+1; i++ {
		endpoint := new(api.Endpoint)
		endpoint.Id = i
		lbImpl.Add(i, endpoint)
	}
	endpoint, err := lbImpl.Pick(nil)
	if err != nil {
		t.Fatal(err)
	}
	if endpoint == nil {
		t.Fatal("pick nil")
	}

	distribMap := make(map[int]int)
	for i := 0; i < 1000000; i++ {
		endpoint, err := lbImpl.Pick(nil)
		if err != nil {
			t.Fatal(err)
		}
		if endpoint == nil {
			t.Fatal("pick nil")
		}
		cnt, ok := distribMap[endpoint.Id]
		if ok {
			distribMap[endpoint.Id] = cnt + 1
		} else {
			distribMap[endpoint.Id] = 1
		}
	}

	t.Log(distribMap)
	cnt := 0
	for range distribMap {
		cnt++
	}
	if cnt < 256+1 {
		t.Fatal("rand is not well distributed")
	}
}
