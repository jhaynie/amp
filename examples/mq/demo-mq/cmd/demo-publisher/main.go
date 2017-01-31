package main

import (
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/examples/mq/demo-mq/messages"
	"github.com/appcelerator/amp/pkg/mq"
	"github.com/appcelerator/amp/pkg/mq/nats-streaming"
	"log"
	"os"
)

var (
	// MQ is the message queuer interface
	MQ mq.Interface
)

const (
	demoQueue = "demo-queue"
)

func main() {
	log.Println(os.Args[0])

	// Connect to message queuer
	log.Println("Connecting to amp MQ")
	MQ = ns.New(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0])
	if err := MQ.Connect(amp.DefaultTimeout); err != nil {
		log.Fatal(err)
	}
	defer MQ.Close()
	log.Println("Connected to amp MQ")

	// Create a demo message
	demoMessage := messages.DemoMessage{
		Data: "test-message",
	}

	// Publish to MQ
	_, err := MQ.PublishAsync(demoQueue, &demoMessage, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Message successfuly published")
}
