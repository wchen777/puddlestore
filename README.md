# Puddlestore README

### Prerequisites for using Tapestry

You will need to use your own Tapestry implementation or the TA implementation of Tapestry for this project (see the handout for more details). However, you will not be using import statements in Go for this. Instead, you will be using a powerful feature known as [the `replace directive` in Go modules](https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/)

Steps:

1. Download the Tapestry implementation to a local folder.
2. Update this line in `go.mod`

```
replace tapestry => /path/to/your/tapestry/implementation/root/folder
```

so that imports of `tapestry` now point to your local folder where you've downloaded Tapestry.

That's it! When running tests in Gradescope, we will automatically rewrite this line to point to our TA implementation, so you will be tested against our implementation and not be penalized for any issues of your own.

### Prerequisites for using Zookeeper

With docker, you can spin up a Zookeeper instance with

```
docker run --rm -p 2181:2181 zookeeper
```

Please refer to the Zookeeper lab for more detailed setup instructions.

### Tests

To run test suite:
`go test ./test -coverpkg ./pkg/... -coverprofile=coverage.out` to test coverage on just `pkg`, then `go tool cover -html=coverage.out` to check coverage.

- `basic_io_test.go` covers basic tests to check if clients can open and close and write to files.
- `list_test.go` tests the client's `List()` function
- `lock_test.go` tests the distributed locks in zookeeper
  - many tests passed individually but were failing when ran with the whole suite, so we commented them out
- `mkdir_test.go` tests the client's `Mkdir()` function (along with edge cases)
- `open_test.go` tests the client's `Open()` function (along with edge cases)
- `read_write_test.go` tests edge cases in reading and writing to files
- `remove_test.go` tests the client's `Remove()` function (along with edge cases)
- `tapestry_test.go` recreates the even distribution and load balancing tests found on Gradescope, and tests overall fault tolerance
  - even distribution tests are failing on Gradescope, so we commented these out
- `no_conn_test.go` tests edge cases regarding having no connections

### Distribution of Work

Will and Mario both worked on the client interface functions and debugged them when necessary, and also both wrote numerous tests.

Will worked on the web client and setting up the gRPC server.

Mario worked on adding replication and load balancing to the tapestry nodes and data blocks.

### Issues with Our Implementation

We are unable to pass the even distribution tests on Gradescope. We are confused on how these are supposed to pass, we load balance such that blocks are distributed evenly across tapestry nodes. When downing a significant number of these nodes, without any replication, then we returned data with "holes" in it, where the data blocks were not able to be read from the downed tapestry nodes. We recreated these tests in `tapestry_test.go`, but we ended up commenting out as they don't pass.

We also have issues with some tests passing individually but not running as a whole. We have ensured that we are cleaning up correctly, but still cause tests to hang or other tests to fail.

### Extra Features

We created a web client.

### How to Run Web Client

1. Switch to `web-client` branch
2. Ensure Zookeeper is running on port 2181
3. Run the server using `go run api/main.go` from the root.
4. Run the following command: `./grpcwebproxy-v0.15.0-osx-x86_64 --backend_addr=:3333 --server_http_debug_port 3334 --allow_all_origins --run_tls_server=false` in the root to ensure the proxy forwards connections from the client to the server.
5. Run the client either by running `npm run start` from within the `puddlestore-web-client` directory, or opening the `index.html` in the `puddlestore-web-client-build` directory.
6. Have fun!

![puddlestore client](img/puddlestore%20sc.png "PuddleStore Client")

`protoc puddlestore.proto \
--js_out=import_style=commonjs,binary:./src/puddlestore \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/puddlestore`
