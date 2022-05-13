package test

import (
	puddlestore "puddlestore/pkg"
	tap "tapestry/pkg"
	"testing"
	"time"
)

// test no zk conn
func TestNoConn(t *testing.T) {
	_, err := puddlestore.CreateCluster(puddlestore.Config{
		ZkAddr:      "localhost:2182", // wrong addr
		BlockSize:   64,
		NumReplicas: 2,
		NumTapestry: 2,
	})

	if err == nil {
		t.Fatal("should have errored")
	}
}

func TestTapClientCreateErr(t *testing.T) {
	_, err := puddlestore.NewTapestry(&tap.Node{}, "invalid string")

	if err == nil {
		t.Fatal("should have errored")
	}
}

// test client exit
func TestClientExit(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	client.Exit()
	time.Sleep(time.Second)
}
