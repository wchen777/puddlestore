package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
)

func TestLockUnit(t *testing.T) {

	zkConn, err := puddlestore.ConnectZk("localhost:2181")
	if err != nil {
		t.Fatal(err)
	}

	_, err = zkConn.Create("/test", []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		t.Fatal(err)
	}

	defer zkConn.Close()

	lock1 := puddlestore.CreateDistLock("/test", zkConn)

	lock2 := puddlestore.CreateDistLock("/test", zkConn)

	errChan := make(chan error)

	go func() {
		err = lock1.Acquire()
		errChan <- err
		time.Sleep(time.Second * 2)
		lock1.Release()
	}()

	go func() {
		err = lock2.Acquire()
		errChan <- err
		time.Sleep(time.Second * 1)
		lock2.Release()
	}()

	err = <-errChan

	if err != nil {
		t.Fatal(err)
	}

	err = <-errChan

	if err != nil {
		t.Fatal(err)
	}

	// clean up
	children, _, _ := zkConn.Children("/")
	for _, c := range children {
		zkConn.Delete("/"+c, -1)
	}
	zkConn.Delete("/", -1)

}

// func TestLock(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	errChan := make(chan error)

// 	go func() {
// 		client, _ := cluster.NewClient()
// 		// if err != nil {
// 		// 	t.Fatal(err)
// 		// }

// 		_, err = client.Open("/test10", true, false)

// 		errChan <- err
// 		client.Exit()
// 	}()

// 	// pull frmo errChan
// 	err = <-errChan
// 	// should not be nil
// 	if err != nil {
// 		t.Fatal("should have err ")
// 	}

// 	go func() {
// 		client2, _ := cluster.NewClient()
// 		_, err = client2.Open("/test5", true, false)
// 		errChan <- err

// 		client2.Exit()
// 	}()

// 	// should not be nil
// 	if err != nil {
// 		t.Fatal("should have err - no create")
// 	}

// }

// func TestBlockWrite(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	client2, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	in := "test"
// 	// error here in write.
// 	if err := writeFile(client, "/a", 0, []byte(in)); err != nil {
// 		t.Fatal(err)
// 	}

// 	in = "out"
// 	// error here in write.
// 	if err := writeFile(client2, "/a", 4, []byte(in)); err != nil {
// 		t.Fatal(err)
// 	}

// 	var out []byte
// 	if out, err = readFile(client2, "/a", 0, 10); err != nil {
// 		t.Fatal(err)
// 	}

// 	output := "testout"
// 	if output != string(out) {
// 		t.Fatalf("Expected: %v, Got: %v", in, string(out))
// 	}

// 	time.Sleep(2 * time.Second)

// 	if out, err = readFile(client2, "/a", 0, 10); err != nil {
// 		t.Fatal(err)
// 	}

// 	if output != string(out) {
// 		t.Fatalf("Expected: %v, Got: %v", in, string(out))
// 	}

// 	in = "testtwo"
// 	// error here in write.
// 	if err := writeFile(client, "/a", 0, []byte(in)); err != nil {
// 		t.Fatal(err)
// 	}

// 	output = "testtwo"
// 	if out, err = readFile(client2, "/a", 0, 10); err != nil {
// 		t.Fatal(err)
// 	}

// 	if output != string(out) {
// 		t.Fatalf("Expected: %v, Got: %v", in, string(out))
// 	}

// 	client.Exit()

// }

// func TestBlock(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = client.Open("/test10", true, true)

// 	// should not be nil
// 	if err != nil {
// 		t.Fatal("should have err ")
// 	}

// 	// should not be nil
// 	if err != nil {
// 		t.Fatal("should not err")
// 	}

// 	in := "one"
// 	err = client.Write(0, 0, []byte(in))

// 	if err != nil {
// 		t.Fatal("should not err")
// 	}

// 	err = client.Close(0)

// 	if err != nil {
// 		t.Fatal("should not err")
// 	}

// 	client.Exit()

// }
