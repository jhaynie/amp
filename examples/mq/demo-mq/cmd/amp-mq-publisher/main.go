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

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// MQ is the message queuer interface
	MQ mq.Interface
)

const (
	demoQueue = "demo-queue"
)

func main() {
	log.Printf("%s (version: %s, build: %s)\n", os.Args[0], Version, Build)

	// Connect to message queuer
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to get hostname: %s", err)
	}
	MQ = ns.New(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0]+"-"+hostname)
	if err := MQ.Connect(amp.DefaultTimeout); err != nil {
		log.Fatal(err)
	}

	// Subscribe to queue
	log.Println("Subscribing to queue:", demoQueue)
	_, err = MQ.Subscribe(demoQueue, messageHandler, &messages.DemoMessage{}, mq.DeliverAllAvailable())
	if err != nil {
		MQ.Close()
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
			MQ.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func messageHandler(msg proto.Message, err error) {
	if err != nil {
		log.Println("Error in message processing:", err)
		return
	}

	demoMessage, ok := msg.(*messages.DemoMessage)
	if !ok {
		log.Println("Error in type assertion")
		return
	}
	log.Println("demoMessage", demoMessage)
}
