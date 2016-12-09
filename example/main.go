package main

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/yanmhlv/service-discovery"
)

func main() {
	serviceDiscovery := Must(NewServiceDiscovery("127.0.0.1:8500"))
	fmt.Println(serviceDiscovery)

	makeServer := func(host string, port int) *http.Server {
		appService := &Service{"foo", host, port}

		serviceDiscovery.Register(appService, "_health", 5*time.Second, 5*time.Second)

		mux := http.NewServeMux()
		mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
			fmt.Println("register", serviceDiscovery.Register(appService, "_health",
				5*time.Second,
				5*time.Second,
			))
		})
		mux.HandleFunc("/deregister", func(w http.ResponseWriter, req *http.Request) {
			fmt.Println("deregister", serviceDiscovery.Deregister(appService))
		})
		mux.HandleFunc("/_health", func(w http.ResponseWriter, req *http.Request) {
			//w.WriteHeader(500)
		})
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("hello"))
			//w.WriteHeader(500)
		})

		return &http.Server{
			Handler: mux,
			Addr:    fmt.Sprintf("%s:%d", host, port),
		}
	}

	done := make(chan bool)
	for i := 3000; i < 3010; i++ {
		go makeServer("127.0.0.1", i).ListenAndServe()
	}
	<-done
}
