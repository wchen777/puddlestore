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

	// create instance of puddlestore server
	ps := puddlestore.PuddleStoreServerInstance{
		Addr: listener.Addr(),
	}
	// initialize the cluster
	ps.InitCluster()

	// create a new grpc server
	server := grpc.NewServer()

	// register the puddlestore server with the grpc server
	puddlestore.RegisterPuddleStoreServer(server, &ps)

	// webhello.RegisterWebHelloServer(server, &node)
	fmt.Println("RPC serving on :3333")
	server.Serve(listener)
}
