package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	client, err := NewConsulClient(address)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Register(name, port)
	if err != nil {
		log.Fatal(err)
	}

	entry, _, err := client.Service(name, "")
	if err != nil {
		log.Fatal(err)
	}

	if len(entry) > 0 {
		services, err := listServices()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(services[name])
	}

	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func listServices() (Services, error) {
	resp, err := http.Get(consulServices)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var services Services
	json.NewDecoder(resp.Body).Decode(&services)

	return services, nil
}

// Service ...
type Service struct {
	ID                string
	Service           string
	Tags              []string
	Address           string
	Port              int
	EnableTagOverride bool
	CreateIndex       int
	ModifyIndex       int
}

// Services ...
type Services map[string]Service

// Client ...
type Client interface {
	// Get a Service from consul
	Service(string, string) ([]string, error)
	// Register a service with local agent
	Register(string, int) error
	// Deregister a service with local agent
	DeRegister(string) error
}

type client struct {
	consul *consul.Client
}

// NewConsulClient returns a Client interface for given consul address
func NewConsulClient(addr string) (*client, error) {
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

// Service return a service
func (c *client) Service(service, tag string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	passingOnly := true
	addrs, meta, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
}
