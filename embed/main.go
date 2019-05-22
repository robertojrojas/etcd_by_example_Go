package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	doneCh := make(chan struct{})
	go startEmbeddedServer(doneCh)

	time.Sleep(10 * time.Second)

	serverURL := "http://localhost:2379"

	clientConfig := clientv3.Config{
		Endpoints: []string{
			serverURL,
		},
		DialTimeout: time.Duration(5 * time.Second),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	client, err := clientv3.New(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Client is connected!!! \n client info: %#v\n", client)

	client.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("CTRL-C to Stop \n")

	for {
		select {
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exciting...", s))
			doneCh <- struct{}{}
			<- doneCh
			os.Exit(0)
		}
	}

}

func startEmbeddedServer(doneCh chan struct{}) {

	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()


	// Handle Server startup notifications
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!\n")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}

 // Handle shutdown notifications
	for {
		select {
		case err := <-e.Err():
			 log.Fatalf("Embedded server failed %v\n", err)
		case <-doneCh:
			fmt.Printf("Done! Stopping Etcd Server!\n")
			e.Server.Stop()
			doneCh <- struct{}{}
		}
	}

}
