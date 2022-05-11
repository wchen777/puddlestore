package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
)

func TestLock(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Open("/test10", true, false)

	// should not be nil
	if err != nil {
		t.Fatal("should have err ")
	}

	client2, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client2.Open("/test5", true, false)

	// should not be nil
	if err != nil {
		t.Fatal("should have err - no create")
	}

}
