package main

import (
	"fmt"
	"net/http"
	"strings"

	rh "github.com/michaelklishin/rabbit-hole/v2"
)

func main() {
	host := ""
	username := ""
	password := ""
	vhost := ""
	queueMessages := ""
	transport := &http.Transport{TLSClientConfig: nil}
	rmqc, err := rh.NewTLSClient(host, username, password, transport)
	if err != nil {
		panic(err)
	}

	_, err = rmqc.PurgeQueue(vhost, queueMessages)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PURGE SUCCESS\n")

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
