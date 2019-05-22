package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
	"context"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <server:port>\n", os.Args[0])
		os.Exit(1)
	}

	endpoints := os.Args[1:]
	fmt.Printf("Endpoints %v\n", endpoints)
	dialTimeout := time.Duration(20 * time.Second)

	clientConfig := clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: dialTimeout,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	cli, err := clientv3.New(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Client is connected!!! \n", cli)

	if _, err = cli.RoleAdd(context.TODO(), "root"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserAdd(context.TODO(), "root", "123"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserGrantRole(context.TODO(), "root", "root"); err != nil {
		log.Fatal(err)
	}

	if _, err = cli.RoleAdd(context.TODO(), "r"); err != nil {
		log.Fatal(err)
	}

	if _, err = cli.RoleGrantPermission(
		context.TODO(),
		"r",   // role name
		"foo", // key
		"zoo", // range end
		clientv3.PermissionType(clientv3.PermReadWrite),
	); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserAdd(context.TODO(), "u", "123"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.UserGrantRole(context.TODO(), "u", "r"); err != nil {
		log.Fatal(err)
	}
	if _, err = cli.AuthEnable(context.TODO()); err != nil {
		log.Fatal(err)
	}

	cliAuth, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		Username:    "u",
		Password:    "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cliAuth.Close()

	if _, err = cliAuth.Put(context.TODO(), "foo1", "bar"); err != nil {
		log.Fatal(err)
	}

	_, err = cliAuth.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.Value("zoo1"), ">", "abc")).
		Then(clientv3.OpPut("zoo1", "XYZ")).
		Else(clientv3.OpPut("zoo1", "ABC")).
		Commit()
	fmt.Println(err)

	// now check the permission with the root account
	rootCli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		Username:    "root",
		Password:    "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rootCli.Close()

	resp, err := rootCli.RoleGet(context.TODO(), "r")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("user u permission: key %q, range end %q\n", resp.Perm[0].Key, resp.Perm[0].RangeEnd)

	if _, err = rootCli.AuthDisable(context.TODO()); err != nil {
		log.Fatal(err)
	}

}
