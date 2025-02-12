package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

func newSimpleServer(address string) *simpleServer {

	serverUrl, err := url.Parse(address)
	handleError(err)

	return &simpleServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
	}

}

type LoadBalancer struct {
	port            int
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port int, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

//simple server methods
func (s *simpleServer) Address() string {
    return s.address
}

func (s *simpleServer) IsAlive() bool {
    _, err := http.Get(s.address)
    return err == nil 
}

func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
    s.proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {

}

func (lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {

}

func main(){
	server1 := newSimpleServer("http://www.google.com")
    server2 := newSimpleServer("https://www.facebook.com")
	server3 := newSimpleServer("http://www.bing.com")

    servers := []Server{server1, server2, server3}

	// Create a new load balancer with servers and start the server
    lb := NewLoadBalancer(8000, servers)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}
    http.HandleFunc("/", handleRedirect)

    fmt.Printf("Starting server on port %d...", lb.port)
    err := http.ListenAndServe(":8000", nil)
    handleError(err)
}