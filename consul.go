package consul

import (
	"context"

	"github.com/hashicorp/consul/api"
	"golang.org/x/sync/errgroup"
)

type (
	Consul struct {
		cl       *api.Client
		Services map[string][]*ConsulService
	}
)

func New(addr string) (*Consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr

	cl, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Consul{cl, make(map[string][]*ConsulService)}, nil
}

func Must(consul *Consul, err error) *Consul {
	if err != nil {
		panic(err)
	}
	return consul
}

func (c *Consul) Register(svc *ConsulService) error {
	serviceID := svc.generateID()

	reg := &api.AgentServiceRegistration{
		Name:    svc.Name,
		ID:      serviceID,
		Address: svc.Host,
		Port:    svc.Port,
		// Check:   svc.Check,
	}

	if err := c.cl.Agent().ServiceRegister(reg); err != nil {
		return err
	}

	id := serviceID
	c.Services[id] = append(c.Services[id], svc)

	return nil
}

func (c *Consul) Deregister() error {
	eg, ctx := errgroup.WithContext(context.Background())
	for name, services := range c.Services {
		_, services := name, services
		eg.Go(func() error {
			for _, service := range services {
				if err := c.cl.Agent().ServiceDeregister(service.generateID()); err != nil {
					return err
				}
			}
			return ctx.Err()
		})
	}

	return eg.Wait()
}
