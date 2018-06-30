package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	consul "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	address        = "http://localhost:8500"
	consulServices = address + "/v1/agent/services"
	name           = "echo"
	port           = 1323
)

func main() {
	client, err := setup()
	if err != nil {
		client.DeRegister(name)
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			client.DeRegister(name)
			os.Exit(1)
		}
	}()

	e := server()
	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func setup() (*client, error) {
	client, err := newConsulClient(address)
	if err != nil {
		return nil, err
	}

	if err = client.Register(name, port); err != nil {
		return nil, err
	}

	services, err := client.consul.Agent().Services()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(services[name])

	return client, nil
}

func server() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
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
	DeRegister(string) error
}

type client struct {
	consul *consul.Client
}

// NewConsulClient returns a Client interface for given consul address
func newConsulClient(addr string) (*client, error) {
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
		ID:   name,
		Name: name,
		Port: port,
	}
	return c.consul.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func (c *client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}
