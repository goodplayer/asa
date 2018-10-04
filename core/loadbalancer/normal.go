package loadbalancer

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/goodplayer/asa/api"
)

const (
	_SUB_ARRAY_SHIFT = 8
	_SUB_ARRAY_MASK  = 2 ^ _SUB_ARRAY_SHIFT - 1
)

type LoadBalancer struct {
	strategy api.PickStrategy

	data    [][]*api.Endpoint
	mapping map[int]int
	cnt     int
	lock    sync.RWMutex
}

func NewLoadBalancer(strategy api.PickStrategy) api.LoadBalancer {
	lb := new(LoadBalancer)
	lb.strategy = strategy
	lb.cnt = 0
	lb.mapping = make(map[int]int)
	//TODO
	return lb
}

func (this *LoadBalancer) Pick(key interface{}) (*api.Endpoint, error) {
	switch this.strategy {
	case api.RANDOM:
		return this.pickRandom()
	case api.ROUND_ROBIN:
		//TODO
	case api.FIXED_KEY:
		//TODO
	}
	panic("implement me")
}

func twoIdx(idx int) (int, int) {
	return idx >> _SUB_ARRAY_SHIFT, idx & _SUB_ARRAY_MASK
}

func (this *LoadBalancer) Add(key interface{}) {
	//TODO
}

func (this *LoadBalancer) Remove(key interface{}) {
	//TODO
}

func (this *LoadBalancer) pickRandom() (*api.Endpoint, error) {
	this.lock.RLock()
	cnt := this.cnt
	if cnt <= 0 {
		this.lock.RUnlock()
		return nil, errors.New("no element found")
	}
	idx := rand.Intn(cnt)
	firstIdx, secondIdx := twoIdx(idx)
	endpoint := this.data[firstIdx][secondIdx]
	this.lock.RUnlock()
	return endpoint, nil
}
