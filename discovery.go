package discovery

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

type (
	// it is connection to consul and storage for services
	ServiceDiscovery struct {
		cl       *api.Client
		Services map[string][]*Service
	}

	// information about service
	Service struct {
		Name    string
		Address string
		Port    int
	}
)

// NewServiceDiscovery creates connect to consul
func NewServiceDiscovery(address string) (*ServiceDiscovery, error) {
	conf := api.DefaultConfig()
	conf.Address = address
	cl, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &ServiceDiscovery{cl, make(map[string][]*Service)}, nil
}

// can panic
func Must(discovery *ServiceDiscovery, err error) *ServiceDiscovery {
	if err != nil {
		panic(err)
	}
	return discovery
}

func (s *Service) hash() string {
	return strings.Join([]string{s.Name, s.Address, fmt.Sprint(s.Port)}, "_")
}

// Register service in consul
func (s *ServiceDiscovery) Register(service *Service, healthPath string,
	interval time.Duration, timeout time.Duration) error {

	u := &url.URL{
		Host:   fmt.Sprintf("%s:%d", service.Address, service.Port),
		Scheme: "http",
		Path:   healthPath,
	}
	if err := s.cl.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      service.hash(),
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Check: &api.AgentServiceCheck{
			Interval: interval.String(),
			Timeout:  timeout.String(),
			HTTP:     u.String(),
		},
	}); err != nil {
		return err
	}

	services := s.Services[service.hash()]
	if services == nil {
		services = make([]*Service, 0)
	}
	s.Services[service.hash()] = append(services, service)

	return nil
}

// Deregister and remove service from consul
func (s *ServiceDiscovery) Deregister(service *Service) error {
	if err := s.cl.Agent().ServiceDeregister(service.hash()); err != nil {
		return err
	}

	s.Services[service.hash()] = make([]*Service, 0)
	return nil
}
