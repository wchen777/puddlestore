package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
)

func TestMkDir(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Mkdir("/test")

	// should not be nil
	if err != nil {
		t.Fatal("should have been made")
	}

	size, _ := client.List("/puddlestore")

	if len(size) != 1 {
		t.Fatal("should only be one directory")
	}
}
