package consul

import (
	"net"
	"strconv"
	"strings"
)

type ConsulService struct {
	Name string
	Host string
	Port int
	Tags []string

	Check *ConsulCheck
}

func NewService(name, host string, port int, tags ...string) *ConsulService {
	return &ConsulService{Name: name, Host: host, Port: port, Tags: tags}
}

func (svc *ConsulService) SetCheck(check ConsulChecker) *ConsulService {
	if check != nil {
		svc.Check = check.BuildCheck()
	}
	return svc
}

func (svc *ConsulService) generateID() string {
	keys := []string{svc.Name, svc.Host, strconv.Itoa(svc.Port)}
	keys = append(keys, svc.Tags...)
	return strings.Join(keys, "|")
}

func (svc *ConsulService) address() string {
	return net.JoinHostPort(svc.Host, strconv.Itoa(svc.Port))
}
