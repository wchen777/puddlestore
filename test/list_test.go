package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
)

func TestList(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Open("/test0", true, false)

	// should not be nil
	if err != nil {
		t.Fatal("should have err - no create")
	}

	_, err = client.Open("/test1", true, false)

	// should not be nil
	if err != nil {
		t.Fatal("should have err - no create")
	}

	lst, err := client.List("/test1")

	if err != nil || len(lst) == 0 {
		t.Fatal("should have err - no list")
	}

}
