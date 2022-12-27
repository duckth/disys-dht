package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	gRPC "github.com/duckth/disys-dht/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var servers []gRPC.HashTableClient

func main() {
	ports := []int64{5000, 5001, 5002}

	for i := 0; i < len(ports); i++ {
		servers = append(servers, ConnectToServer(ports[i]))
	}

	readInput()
}

func readInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := strings.Split(scanner.Text(), " ")

		if in[0] == "put" {
			if len(in) != 3 {
				log.Printf("Input needs to be of form: put <key> <value>")
				continue
			}
			key, err1 := strconv.ParseInt(in[1], 10, 64)
			val, err2 := strconv.ParseInt(in[2], 10, 64)

			if err1 != nil {
				log.Printf("error: %v", err1)
			} else if err2 != nil {
				log.Printf("error: %v", err2)
			} else {
				Put(servers, key, val)
			}
		} else if in[0] == "get" {
			if len(in) != 2 {
				log.Printf("Input needs to be of form: get <key>")
				continue
			}
			key, err1 := strconv.ParseInt(in[1], 10, 64)

			if err1 != nil {
				log.Printf("error: %v", err1)
			} else {
				Get(servers, key)
			}
		} else {
			log.Printf("Input needs to be one of the following:")
			log.Printf("put <key> <value>")
			log.Printf("get <key>")
		}
	}
}

func Put(servers []gRPC.HashTableClient, key int64, value int64) {
	for i := 0; i < len(servers); i++ {
		server := servers[i]

		response, _ := server.Put(context.Background(), &gRPC.PutRequest{Key: key, Value: value})

		if response == nil {
			log.Printf("Received no response from server %d", i)
		} else if !response.Success {
			log.Printf("Unsuccessful put request of { %d => %d } for server %d", key, value, i)
		} else {
			log.Printf("Successful put request of { %d => %d } for server %d", key, value, i)
		}
	}
}

func Get(servers []gRPC.HashTableClient, key int64) {
	hasReceivedResponse := false

	for i := 0; i < len(servers); i++ {
		server := servers[i]

		if hasReceivedResponse {
			continue
		}

		response, _ := server.Get(context.Background(), &gRPC.GetRequest{Key: key})

		if response == nil {
			log.Printf("Received no response from server %d", i)
			continue
		}

		log.Printf("GET: { %d => %d }", key, response.Value)
		hasReceivedResponse = true
	}
}

func ConnectToServer(port int64) gRPC.HashTableClient {
	var opts []grpc.DialOption
	var target = fmt.Sprintf("localhost:%d", port)

	opts = append(
		opts,
		grpc.WithBlock(),
		grpc.WithTimeout(1*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	fmt.Printf("Dialing on %s \n", target)
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("Fail to dial: %v\n", err)
	}

	return gRPC.NewHashTableClient(conn)
}
