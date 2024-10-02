package mq

import "github.com/streadway/amqp"

// amqp://user:pass@host:10000/vhost
type Cfg struct {
	MaxAttempts      int
	User             string
	Password         string
	Host             string
	Port             string
	Vhost            string
	ExchangeSettigns ExchangeSettigns
}

type ExchangeSettigns struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func DefaultExchangeSettings(name, kind string) ExchangeSettigns {
	return ExchangeSettigns{
		Name:       name,
		Kind:       kind,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}
