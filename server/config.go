package server

import newrelictracing "github.com/anant-sharma/go-utils/new-relic/tracing"

type Config struct {
	Grpc struct {
		Host string
		Port int
	}
	HTTP struct {
		Host string
		Port int
	}
	NewRelic newrelictracing.Config
}
