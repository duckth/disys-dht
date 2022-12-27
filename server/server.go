package main

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"

	gRPC "github.com/duckth/disys-dht/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedHashTableServer
	port      string
	hashTable map[int64]int64
	mutex     sync.Mutex
}

var port = flag.String("port", "5000", "Server port")

func main() {
	flag.Parse()
	listen, _ := net.Listen("tcp", "localhost:"+*port)
	grpcServer := grpc.NewServer()
	hashTableServer := &Server{
		port:      *port,
		hashTable: make(map[int64]int64),
	}

	gRPC.RegisterHashTableServer(grpcServer, hashTableServer)
	log.Printf("Listening on port %s...", *port)
	grpcServer.Serve(listen)
}

func (s *Server) Put(ctx context.Context, req *gRPC.PutRequest) (*gRPC.PutResponse, error) {
	log.Printf("Receiving put request: { %d => %d }", req.Key, req.Value)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.hashTable[req.Key] = req.Value
	return &gRPC.PutResponse{Success: true}, nil
}

func (s *Server) Get(ctx context.Context, req *gRPC.GetRequest) (*gRPC.GetResponse, error) {
	log.Printf("Receiving get request of key %d", req.Key)

	value, present := s.hashTable[req.Key]

	if !present {
		value = 0
	}

	return &gRPC.GetResponse{Value: value}, nil
}
