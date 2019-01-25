package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
)

type (
	ConsulCheck = api.AgentServiceCheck

	ConsulChecker interface {
		BuildCheck() *ConsulCheck
	}

	GRPCheck struct {
		Addr     string
		Interval time.Duration
		Timeout  time.Duration
	}

	HTTPCheck struct {
		URL      string
		Interval time.Duration
		Timeout  time.Duration
	}
)

var (
	_ ConsulChecker = &GRPCheck{}
	_ ConsulChecker = &HTTPCheck{}
)

func (c *GRPCheck) BuildCheck() *ConsulCheck {
	return &ConsulCheck{
		GRPCUseTLS: false,
		GRPC:       c.Addr,
		Timeout:    c.Timeout.String(),
		Interval:   c.Interval.String(),
	}
}

func (c *HTTPCheck) BuildCheck() *ConsulCheck {
	return &ConsulCheck{
		Timeout:  c.Timeout.String(),
		Interval: c.Interval.String(),
		HTTP:     c.URL,
	}
}
