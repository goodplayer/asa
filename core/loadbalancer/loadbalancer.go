package loadbalancer

import (
	"github.com/goodplayer/asa/api"
)

var _ api.LoadBalancer = new(LoadBalancer)
var _ api.LoadBalancer = new(WeightLoadBalancer)

func NewBalancer(strategy api.PickStrategy) api.LoadBalancer {
	switch strategy {
	case api.RANDOM:
		return new(LoadBalancer)
	}
	panic(api.ErrUnsupportedStrategy)
}

type WeightLoadBalancer struct {
}

func (this *WeightLoadBalancer) Pick(key interface{}) (*api.Endpoint, error) {
	panic("implement me")
}
