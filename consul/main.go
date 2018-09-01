package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	consul "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	address = "http://localhost:8500"
	name    = "echo"
	port    = 1323
)

func main() {
	// Get a new client
	client, err := newClient(address)
	if err != nil {
		panic(err)
	}

	// register echo server to the services
	if err = client.Register(name, port); err != nil {
		panic(err)
	}

	// get services
	services, err := client.consul.Agent().Services()
	if err != nil {
		panic(err)
	}
	s := services[name]
	fmt.Printf("Service '%s' running on '%s:%d'\n", name, s.Address, s.Port)

	// Key Value API
	kv := client.consul.KV()
	// put a key value pair
	p := &consul.KVPair{Key: "max_connections", Value: []byte("100")}
	if _, err = kv.Put(p, nil); err != nil {
		panic(err)
	}
	// get the key value pair
	pair, _, err := kv.Get("max_connections", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v: %s\n", pair.Key, pair.Value)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			client.Deregister(name)
			os.Exit(1)
		}
	}()

	e := server()
	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func server() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	return e
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Client ...
type Client interface {
	// Register a service with local agent
	Register(string, int) error
	// Deregister a service with local agent
	Deregister(string) error
}

type client struct {
	consul *consul.Client
}

// NewClient returns a Client interface for given consul address
func newClient(addr string) (*client, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	c, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &client{consul: c}, nil
}

// Register a service with consul local agent
func (c *client) Register(name string, port int) error {
	reg := &consul.AgentServiceRegistration{
		ID:      name,
		Name:    name,
		Port:    port,
		Address: "localhost",
	}
	return c.consul.Agent().ServiceRegister(reg)
}

// Deregister a service with consul local agent
func (c *client) Deregister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}
