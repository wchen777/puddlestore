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

	size, _ := client.List("")

	// should get everything below /puddlestore
	if len(size) != 1 {
		t.Fatal("should only be one directory")
	}
}

func TestMkDirInvalid(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Mkdir("/dev/test")

	// should not be able to create in nonexistent directory
	if err == nil {
		t.Fatal("dir should not have been made")
	}

	// create a file
	fd, err := client.Open("/notadir", true, false)
	if err != nil {
		t.Fatal("file should have been made")
	}

	// close the file
	err = client.Close(fd)
	if err != nil {
		t.Fatal("file should have been closed")
	}

	// should not be able to create under a file
	err = client.Mkdir("/notadir/test")
	if err == nil {
		t.Fatal("dir should not have been made")
	}

}

func TestMkDirDuplicate(t *testing.T) {
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
		t.Fatal("dir should have been made")
	}

	err = client.Mkdir("/test")

	// should error
	if err == nil {
		t.Fatal("dir should not have been made")
	}

}
