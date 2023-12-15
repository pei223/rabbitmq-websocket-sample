package main

type Config struct {
	RabbitMqURL string `envconfig:"RABBITMQ_URL" default:"rabbitmquser:rabbitmqpass@localhost:5672"`
}
