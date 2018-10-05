package loadbalancer

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/goodplayer/asa/api"
)

const (
	_SUB_ARRAY_SHIFT = 8
	_SUB_ARRAY_MASK  = (1 << _SUB_ARRAY_SHIFT) - 1
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
	//TODO dynamic increase
	lb.data = make([][]*api.Endpoint, 1024)
	for i := range lb.data {
		lb.data[i] = make([]*api.Endpoint, 1<<_SUB_ARRAY_SHIFT)
	}
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

func (this *LoadBalancer) Add(key interface{}, node *api.Endpoint) {
	switch this.strategy {
	case api.RANDOM:
		this.addRandom(key.(int), node)
		return
	case api.ROUND_ROBIN:
		//TODO
		return
	case api.FIXED_KEY:
		//TODO
		return
	}
	panic("implement me")
}

func (this *LoadBalancer) Remove(key interface{}) {
	switch this.strategy {
	case api.RANDOM:
		this.removeRandom(key.(int))
		return
	case api.ROUND_ROBIN:
		//TODO
		return
	case api.FIXED_KEY:
		//TODO
		return
	}
	panic("implement me")
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

func (this *LoadBalancer) addRandom(i int, endpoint *api.Endpoint) {
	this.lock.Lock()
	cnt := this.cnt
	firstIdx, secondIdx := twoIdx(cnt)
	this.data[firstIdx][secondIdx] = endpoint
	this.cnt++
	this.mapping[i] = cnt
	this.lock.Unlock()
}

func (this *LoadBalancer) removeRandom(i int) {
	this.lock.Lock()
	idx, ok := this.mapping[i]
	if !ok {
		this.lock.Unlock()
		return
	}

	// last
	firstIdx, secondIdx := twoIdx(this.cnt - 1)
	last := this.data[firstIdx][secondIdx]

	newFirstIdx, newSecondIdx := twoIdx(idx)
	this.data[newFirstIdx][newSecondIdx] = last
	this.lock.Unlock()
}
