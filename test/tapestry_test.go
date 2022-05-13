package test

// func TestKill1in3(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 2,
// 		NumTapestry: 3,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttesttesttesttesttesttesttesttetst")) // spans multiple blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill one node
// 	// get random num from 0 to 4
// 	cluster.GetTapestryNodes()[1].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 16)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtesttesttest" {
// 		t.Fatalf("Expected: %v, Got: %v", "testtesttesttest", string(data))
// 	}

// }

// func TestKill1in5(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 2,
// 		NumTapestry: 5,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttesttesttesttesttesttesttesttetst")) // spans multiple blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill one node
// 	// get random num from 0 to 4
// 	cluster.GetTapestryNodes()[2].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 16)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtesttesttest" {
// 		t.Fatalf("Expected: %v, Got: %v", "testtesttesttest", string(data))
// 	}

// 	err = client.Close(fd1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// }

// func TestKill2in5(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 5,
// 		NumTapestry: 5,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttesttesttesttesttesttesttesttetst")) // spans multiple blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill two nodes
// 	cluster.GetTapestryNodes()[2].GracefulExit()
// 	cluster.GetTapestryNodes()[4].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 16)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtesttesttest" {
// 		t.Fatalf("Expected: %v, Got: %v\n", "testtesttesttest", string(data))
// 	}

// }

// func TestKill2in5OneReplica(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 1,
// 		NumTapestry: 5,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttesttesttesttesttesttesttesttetst")) // spans multiple blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill two nodes
// 	cluster.GetTapestryNodes()[0].GracefulExit()
// 	cluster.GetTapestryNodes()[3].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 16)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtesttesttest" {
// 		t.Fatalf("Expected: %v, Got: %v\n", "testtesttesttest", string(data))
// 	}

// }

// func TestKill2in5LowBlocks(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 2,
// 		NumTapestry: 5,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttest")) // spans few blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill two nodes
// 	cluster.GetTapestryNodes()[0].GracefulExit()
// 	cluster.GetTapestryNodes()[1].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 8)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtest" {
// 		t.Fatalf("Expected: %v, Got: %v\n", "testtest", string(data))
// 	}
// }

// func TestKill3in5OneReplica(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 1,
// 		NumTapestry: 5,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttesttesttesttesttesttesttesttetst")) // spans multiple blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill three nodes
// 	cluster.GetTapestryNodes()[0].GracefulExit()
// 	cluster.GetTapestryNodes()[1].GracefulExit()
// 	cluster.GetTapestryNodes()[2].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 16)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtesttesttest" {
// 		t.Fatalf("Expected: %v, Got: %v\n", "testtesttesttest", string(data))
// 	}

// }

// func TestKill3in5(t *testing.T) {
// 	cluster, err := puddlestore.CreateCluster(puddlestore.Config{
// 		BlockSize:   8,
// 		NumReplicas: 2,
// 		NumTapestry: 5,
// 		ZkAddr:      "localhost:2181", // restore to localhost:2181 before submitting
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer cluster.Shutdown()

// 	client, err := cluster.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// open a test files
// 	fd0, err := client.Open("/test0", true, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// write to the file
// 	err = client.Write(fd0, 0, []byte("testtesttesttesttesttesttesttesttesttesttetst")) // spans multiple blocks

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the file
// 	err = client.Close(fd0)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// kill three nodes
// 	cluster.GetTapestryNodes()[0].GracefulExit()
// 	cluster.GetTapestryNodes()[1].GracefulExit()
// 	cluster.GetTapestryNodes()[2].GracefulExit()

// 	// reopen the file
// 	fd1, err := client.Open("/test0", false, true)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// read the file
// 	data, err := client.Read(fd1, 0, 16)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(data) != "testtesttesttest" {
// 		t.Fatalf("Expected: %v, Got: %v\n", "testtesttesttest", string(data))
// 	}

// }
