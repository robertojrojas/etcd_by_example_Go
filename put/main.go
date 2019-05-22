package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <server:port>\n", os.Args[0])
		os.Exit(1)
	}

	serverURL := os.Args[1]

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

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	putResponse, err := client.Put(ctx, "golang", "philly")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Put Response: %#v\n", putResponse)

}
