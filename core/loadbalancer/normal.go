package loadbalancer

import "github.com/goodplayer/asa/api"

type LoadBalancer struct {
	strategy api.PickStrategy
}

func (this *LoadBalancer) Pick(key interface{}) (*api.Endpoint, error) {
	switch this.strategy {
	case api.RANDOM:
		//TODO
	case api.ROUND_ROBIN:
		//TODO
	case api.FIXED_KEY:
		//TODO
	}
	panic("implement me")
}
