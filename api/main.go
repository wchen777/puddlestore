package main

import (
	"fmt"
	"net"
	puddlestore "puddlestore/pkg"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		panic(err)
	}
	// node := webhello.MyWebHelloServer{Addr: listener.Addr()}
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	server := grpc.NewServer()

	// webhello.RegisterWebHelloServer(server, &node)
	fmt.Println("RPC serving on :3333")
	server.Serve(listener)
}
