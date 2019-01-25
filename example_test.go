package consul

import "time"

func ExampleGRPCCheck() {
	svc := NewService("greeter", "localhost", 1234, "hello", "world")
	svc.SetCheck(&GRPCheck{"localhost:1234", time.Second, time.Second})

	consul := Must(New("localhost:8500"))
	consul.Register(svc)
	defer consul.Deregister()
}

func ExampleHTTPCheck() {
	svc := NewService("greeter", "localhost", 1234, "hello", "world")
	svc.SetCheck(&HTTPCheck{"http://localhost:1234/health", time.Second, time.Second})

	consul := Must(New("localhost:8500"))
	consul.Register(svc)
	defer consul.Deregister()
}