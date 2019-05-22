package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
	"google.golang.org/grpc"
)

func main() {

	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go startEmbeddedEtcd(doneCh, &wg)
	wg.Wait() // wait for server to start

	go runClient()

	handleStopRequest(doneCh)

}

func runClient() {
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
	defer client.Close()
	fmt.Printf("Client is connected!!! \n client info: %#v\n", client)
}

func startEmbeddedEtcd(doneCh chan struct{}, wg *sync.WaitGroup) {

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
		log.Fatalf("Server took too long to start!")
	}

	wg.Done()

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

func handleStopRequest(doneCh chan struct{}) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("CTRL-C to Stop \n")

	for {
		select {
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exciting...", s))
			doneCh <- struct{}{}
			<-doneCh
			os.Exit(0)
		}
	}
}
