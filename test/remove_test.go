package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
)

func TestRemoveNonexist(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Remove("no exist")

	// should not be nil
	if err == nil {
		t.Fatal("should have gotten error.")
	}
}

func TestRemoveFile(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	fd, err := client.Open("/test00", true, true)

	if err != nil {
		t.Fatal("should not have error")
	}

	err = client.Close(fd)

	if err != nil {
		t.Fatal("should not have error")
	}

	// should be able to remove file.
	err = client.Remove("/test00")

	// should not be nil
	if err != nil {
		t.Fatal(err)
	}

	// does this file still exist? it should not.

	_, err = client.List("/test00")

	if err == nil {
		t.Fatal("should have returned error since doesn't exist anymore.")
	}
}
