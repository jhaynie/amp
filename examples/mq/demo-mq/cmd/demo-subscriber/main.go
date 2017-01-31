package main

import (
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/examples/mq/demo-mq/messages"
	"github.com/appcelerator/amp/pkg/mq"
	"github.com/appcelerator/amp/pkg/mq/nats-streaming"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"os/signal"
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

	// Connect to MQ
	log.Println("Connecting to amp MQ")
	MQ = ns.New(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0])
	if err := MQ.Connect(amp.DefaultTimeout); err != nil {
		log.Fatal(err)
	}
	defer MQ.Close()
	log.Println("Connected to amp MQ")

	// Subscribe to queue
	log.Println("Subscribing to queue:", demoQueue)
	_, err := MQ.Subscribe(demoQueue, messageHandler, &messages.DemoMessage{}, mq.DeliverAllAvailable())
	if err != nil {
		log.Fatalln("Unable to subscribe to queue", err)
	}
	log.Println("Subscribed to queue:", demoQueue)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Println("\nReceived an interrupt, unsubscribing and closing connection...")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func messageHandler(msg proto.Message, err error) {
	if err != nil {
		log.Println("Error in message processing", err)
		return
	}
	demoMessage, ok := msg.(*messages.DemoMessage)
	if !ok {
		log.Println("Error in type assertion")
		return
	}
	log.Println("demoMessage", demoMessage)
}
