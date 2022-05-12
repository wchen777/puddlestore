package test

import (
	puddlestore "puddlestore/pkg"
	"testing"
	"time"
)

func writeFile(client puddlestore.Client, path string, offset uint64, data []byte) error {
	fd, err := client.Open(path, true, true)
	if err != nil {
		return err
	}
	defer client.Close(fd)

	return client.Write(fd, offset, data)
}

func readFile(client puddlestore.Client, path string, offset, size uint64) ([]byte, error) {
	fd, err := client.Open(path, true, false)
	if err != nil {
		return nil, err
	}
	defer client.Close(fd)

	return client.Read(fd, offset, size)
}

func TestReadWrite(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	in := "test"
	// error here in write.
	if err := writeFile(client, "/b", 0, []byte(in)); err != nil {
		t.Fatal(err)
	}

	var out []byte
	if out, err = readFile(client, "/b", 0, 5); err != nil {
		t.Fatal(err)
	}

	if in != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

	client2, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	in = "one"
	// error here in write.
	if err := writeFile(client2, "/b", 4, []byte(in)); err != nil {
		t.Fatal(err)
	}

	if out, err = readFile(client2, "/b", 0, 10); err != nil {
		t.Fatal(err)
	}

	output := "testone"
	if output != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

	client3, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	in = "two"
	// error here in write.
	if err := writeFile(client3, "/b", 7, []byte(in)); err != nil {
		t.Fatal(err)
	}

	if out, err = readFile(client3, "/b", 0, 14); err != nil {
		t.Fatal(err)
	}

	output = "testonetwo"
	if output != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}
}

func TestFillBlock(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// 64 bytes
	in := "1111111111111111111111111111111111111111111111111111111111111111"
	// error here in write.
	if err := writeFile(client, "/a", 0, []byte(in)); err != nil {
		t.Fatal(err)
	}

	var out []byte
	if out, err = readFile(client, "/a", 0, 65); err != nil {
		t.Fatal(err)
	}

	if in != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

	// two blocks(32 bytes)
	in = "11111111111111111111111111111111"

	// error here in write.
	if err := writeFile(client, "/a", 64, []byte(in)); err != nil {
		t.Fatal(err)
	}

	if out, err = readFile(client, "/a", 20, 100); err != nil {
		t.Fatal(err)
	}

	exp_out := "1111111111111111111111111111111111111111111111111111111111111111111111111111"
	if exp_out != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

}

func TestLargeBlock(t *testing.T) {
	cluster, err := puddlestore.CreateCluster(puddlestore.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Shutdown()

	client, err := cluster.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// 64 bytes
	in := "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
	// error here in write.
	if err := writeFile(client, "/a", 0, []byte(in)); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	var out []byte
	if out, err = readFile(client, "/a", 0, 999); err != nil {
		t.Fatal(err)
	}

	if in != string(out) {
		t.Fatalf("Expected: %v, Got: %v", in, string(out))
	}

}
