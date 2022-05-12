package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
)

// - read empty file
func TestReadEmptyFile(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	fd, err := client.Open("/test.txt", true, false)

	if err != nil {
		t.Fatal("should not have err on open")
	}

	b, err := client.Read(fd, 0, 10)

	if err != nil {
		t.Fatal("should have err on read")
	}

	if len(b) != 0 {
		t.Fatal("should have 0 bytes read")
	}

	client.Close(fd)

}

// - read non open file (random fd), - write to non open file (random fd)
func TestReadWriteNonOpenFile(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Read(5, 2, 10)

	if err == nil {
		t.Fatal("should have err on read")
	}

	err = client.Write(5, 0, []byte("test"))

	if err == nil {
		t.Fatal("should have err on write")
	}

}

// - write to non dirty file
func TestWriteNonDirtyFile(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	fd, err := client.Open("/test.txt", true, false)

	if err != nil {
		t.Fatal("should not have err on open")
	}

	err = client.Write(fd, 0, []byte("test"))

	if err == nil {
		t.Fatal("should have err on write")
	}

}
