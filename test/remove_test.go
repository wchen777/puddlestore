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

func TestRemoveDir(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())

	if err != nil {
		t.Fatal(err)
	}

	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Mkdir("/thisisadir")

	if err != nil {
		t.Fatal("should not have errored on mkdir")
	}

	err = client.Remove("/thisisadir")

	// should not error out
	if err != nil {
		t.Fatal("should have gotten error.")
	}

	// list should return nothing

	_, err = client.List("/thisisadir")

	if err == nil {
		t.Fatal("should have returned error since doesn't exist anymore.")
	}

	child, err := client.List("/")

	if err != nil {
		t.Fatal("should not have gotten error.")
	}

	if len(child) != 0 {
		t.Fatal("should have returned no children.")
	}
}

// remove directory with children
func TestRemoveDirWithChildren(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())

	if err != nil {
		t.Fatal(err)
	}

	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Mkdir("/thisisadir")

	if err != nil {
		t.Fatal("should not have errored on mkdir")
	}

	err = client.Mkdir("/thisisadir/child1")

	if err != nil {
		t.Fatal("should not have errored on mkdir")
	}

	err = client.Mkdir("/thisisadir/child2")

	if err != nil {
		t.Fatal("should not have errored on mkdir")
	}

	fd, err := client.Open("/thisisadir/child1/a.txt", true, true)

	if err != nil {
		t.Fatal("should not have errored on open")
	}

	err = client.Close(fd)

	if err != nil {
		t.Fatal("should not have errored on close")
	}

	fd, err = client.Open("/thisisadir/childfile.txt", true, false)

	if err != nil {
		t.Fatal("should not have errored on open")
	}

	err = client.Close(fd)

	if err != nil {
		t.Fatal("should not have errored on close")
	}

	err = client.Remove("/thisisadir")

	// should not error out
	if err != nil {
		t.Fatal("should not have gotten error: ", err)
	}

}
