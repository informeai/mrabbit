package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	rh "github.com/michaelklishin/rabbit-hole/v2"
)

var helper = `
NAME:
  mrabbit
DESCRIPTION:
  utility for manager rabbit mq.
COMMANDS:
  purge-queue -> purge messages the queue.
                 Ex: mrabbit -h [HOST] -u [USERNAME] -p [PASSWORD] purge-queue [VHOST] [NAME THE QUEUE]
  list-queues -> list queues.
                 Ex: mrabbit -h [HOST] -u [USERNAME] -p [PASSWORD] list-queues
  delete-queues-unused -> delete queues usused.
                 Ex: mrabbit -h [HOST] -u [USERNAME] -p [PASSWORD] delete-queues-unused

USAGE:
  mrabbit -h [HOST] -u [USERNAME] -p [PASSWORD] [COMMAND] [OPTIONS]
`

func main() {
	var host string
	var username string
	var password string

	flag.StringVar(&host, "h", "", "host the rabbit")
	flag.StringVar(&username, "u", "", "username the rabbit")
	flag.StringVar(&password, "p", "", "password the rabbit")

	flag.Parse()
	options := []string{host, username, password}
	for _, opt := range options {
		if len(opt) == 0 {
			fmt.Printf("%s\n", helper)
			os.Exit(0)
		}
	}
	if len(os.Args) < 8 {
		fmt.Printf("%s\n", helper)
		os.Exit(0)
	}
	transport := &http.Transport{TLSClientConfig: nil}
	rmqc, err := rh.NewTLSClient(host, username, password, transport)
	if err != nil {
		panic(err)
	}

	if os.Args[7] == "purge-queue" {
		_, err := rmqc.PurgeQueue(os.Args[8], os.Args[9])
		if err != nil {
			panic(err)
		}
		fmt.Printf("PURGE SUCCESS\n")

	} else if os.Args[7] == "list-queues" {
		queues, err := rmqc.ListQueues()
		if err != nil {
			panic(err)
		}
		for _, queue := range queues {
			fmt.Printf("Vhost: %s Name: %s\n", queue.Vhost, queue.Name)
		}
	} else if os.Args[7] == "delete-queues-unused" {
		queues, err := rmqc.ListQueues()
		if err != nil {
			panic(err)
		}
		for _, queue := range queues {
			if strings.Contains(queue.Name, "amq.gen") {
				fmt.Printf("Vhost: %s Name: %s\n", queue.Vhost, queue.Name)
				_, err := rmqc.DeleteQueue(queue.Vhost, queue.Name)
				if err != nil {
					fmt.Printf("error: %s\n", err.Error())
				}
				fmt.Printf("SUCESS REMOVE QUEUE: %s\n", queue.Name)
			}
		}
	}

}
