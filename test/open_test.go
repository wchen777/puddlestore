package test

import (
	"fmt"
	puddlestore "puddlestore/pkg"
	"testing"
	"time"
)

// - open a file with create under a file
func TestOpenFileCreateUnderFile(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	defer cluster.Shutdown()

	client, err := cluster.NewClient()

	if err != nil {
		t.Fatal(err)
	}

	fd, err := client.Open("/afile.txt", true, false)

	if err != nil {
		t.Fatal("should not have err on open")
	}

	err = client.Close(fd)

	if err != nil {
		t.Fatal("should not have err on close")
	}

	time.Sleep(1 * time.Second)

	_, err = client.Open("/afile.txt/anotherfile.txt", true, false)

	if err == nil {
		t.Fatal("should have err on open")
	}

}

// - open non existent file without create

func TestOpenFileNonExistentNoCreate(t *testing.T) {

	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	defer cluster.Shutdown()

	client, err := cluster.NewClient()

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Open("/noexist.txt", false, false)

	if err == nil {
		t.Fatal("should not be able to create file")
	}

}

// - no mem
func TestOpenFileNoMem(t *testing.T) {

	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	defer cluster.Shutdown()

	maxFiles := 1

	client, err := cluster.NewClientTest(maxFiles)

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < maxFiles; i++ {
		f, err := client.Open("/file"+fmt.Sprint(i)+".c", true, false)
		if err != nil {
			t.Fatal("should not be able to create file")
		}
		err = client.Close(f)
		if err != nil {
			t.Fatal("should not have err on close")
		}

	}

	// open a new file
	_, err = client.Open("/a.txt", true, false)

	// should err
	if err == nil {
		t.Fatal("should not be able to create file due to no mem")
	}
}

// - open existing file that is a dir

func TestOpenFileExistingDir(t *testing.T) {
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

	if err != nil {
		t.Fatal("should not have err on mkdir")
	}

	_, err = client.Open("/test", false, false)

	if err == nil {
		t.Fatal("should not be able to open dir")
	}

}
