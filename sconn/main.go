package main

import (
   "fmt"
	 "time"
	 "log"
	 "os"
	 "go.etcd.io/etcd/clientv3"
	 "google.golang.org/grpc"
	 "crypto/tls"
	"crypto/x509"
	 "github.com/lizrice/secure-connections/utils"
	 "io/ioutil"
)

func main() {

  if len(os.Args) < 5 {
	   fmt.Printf("usage: %s <server:port> <ca-cert> <client-cert>\n", os.Args[0])
		 os.Exit(1)
	}

	serverURL := os.Args[1]
	caCertPath := os.Args[2]
	clientCertPath := os.Args[3]
	clientKeyPath := os.Args[4]

	cp, _ := x509.SystemCertPool()
	data, _ := ioutil.ReadFile(caCertPath)
	cp.AppendCertsFromPEM(data)
	tlsConfig := &tls.Config {
	  RootCAs: cp,
		GetClientCertificate: utils.ClientCertReqFunc(clientCertPath, clientKeyPath),
		VerifyPeerCertificate: utils.CertificateChains,
	}

	clientConfig := clientv3.Config{
	    Endpoints: []string {
			   serverURL,
			},
			DialTimeout:  time.Duration(5 * time.Second),
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
			TLS: tlsConfig,
	}
	client, err := clientv3.New(clientConfig)
	if err != nil {
	   log.Fatal(err)
	}

	fmt.Printf("Client is connected!!! \n client info: %#v\n", client)


	client.Close()

}
