package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
	"time"
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

func TestBlockWrite(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	client2, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	in := "test"
	// error here in write.
	if err := writeFile(client, "/a", 0, []byte(in)); err != nil {
		t.Fatal(err)
	}

	in = "out"
	// error here in write.
	if err := writeFile(client2, "/a", 4, []byte(in)); err != nil {
		t.Fatal(err)
	}

	var out []byte
	if out, err = readFile(client2, "/a", 0, 10); err != nil {
		t.Fatal(err)
	}

	output := "testout"
	if output != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

	time.Sleep(2 * time.Second)

	if out, err = readFile(client2, "/a", 0, 10); err != nil {
		t.Fatal(err)
	}

	if output != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

	in = "testtwo"
	// error here in write.
	if err := writeFile(client, "/a", 0, []byte(in)); err != nil {
		t.Fatal(err)
	}

	output = "testtwo"
	if out, err = readFile(client2, "/a", 0, 10); err != nil {
		t.Fatal(err)
	}

	if output != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

}

func TestBlock(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Open("/test10", true, true)

	// should not be nil
	if err != nil {
		t.Fatal("should have err ")
	}

	// client2, err := cluster.NewClient()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// err := go func() {
	// 	client2.Open("/test10", true, false)
	// 	var out []byte
	// 	out, err = client2.Read(0, 0, 5)

	// 	if err != nil {
	// 		t.Fatal("should not err")
	// 	}

	// 	if "one" != string(out) {
	// 		t.Fatalf("Expected: %v, Got: %v", "one", string(out))
	// 	}

	// 	return nil
	// }()

	// should not be nil
	if err != nil {
		t.Fatal("should not err")
	}

	in := "one"
	err = client.Write(0, 0, []byte(in))

	if err != nil {
		t.Fatal("should not err")
	}

	err = client.Close(0)

	if err != nil {
		t.Fatal("should not err")
	}

}
