package api

import "errors"

type PickStrategy int

const (
	RANDOM PickStrategy = iota
	ROUND_ROBIN
	FIXED_KEY
)

var (
	ErrUnsupportedStrategy = errors.New("unsupported strategy")
)

type LoadBalancer interface {
	Pick(key interface{}) (*Endpoint, error)
}

type Endpoint struct {
	Id int
}
