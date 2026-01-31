package cli

import (
	"errors"
	"rsshub/internal/domain"
	"rsshub/internal/service"
	"strconv"
)

type Cli struct {
	service *service.Service
	args    []string
}

func GetCli(service *service.Service, args []string) *Cli {
	return &Cli{
		service: service,
		args:    args,
	}
}

func (c *Cli) SetInterval(newInterval string) error {
	return c.service.SetInterval(newInterval)
}

func (c *Cli) SetWorkerNs(numWorkers int) error {
	return c.service.SetWorkerNs(numWorkers)
}

func (c *Cli) AddNewNameAndUrl(NewUrl domain.NameAndUrl) error {
	return c.service.AddNewNameAndUrl(NewUrl)
}

func (c *Cli) Start() error {
	return c.service.Start()
}

func (c *Cli) Run() error {
	args := c.args
	if len(args) == 1 && args[0] == "fetch" {
		return c.Start()
	}
	if len(args) == 5 && args[0] == "add" && args[1] == "--name" && args[3] == "--url" {
		var newurl domain.NameAndUrl
		newurl.Name = args[2]
		newurl.Url = args[4]
		return c.AddNewNameAndUrl(newurl)
	}
	if len(args) == 2 && args[0] == "set-interval" {
		return c.SetInterval(args[1])
	}
	if len(args) == 2 && args[0] == "set-workers" {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		return c.SetWorkerNs(n)
	}
	return errors.New("Wrong Command")
}
