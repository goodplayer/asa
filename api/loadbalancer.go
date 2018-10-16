package api

import "errors"

type PickStrategy int

const (
	RANDOM      PickStrategy = iota
	ROUND_ROBIN
	FIXED_KEY
)

var (
	ErrUnsupportedStrategy = errors.New("unsupported strategy")
)

type LoadBalancer interface {
	Pick(key interface{}) (*Endpoint, error)
	Add(key interface{}, node *Endpoint)
	Remove(key interface{})
}

type Endpoint struct {
	Id int64
}
