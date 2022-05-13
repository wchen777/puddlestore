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


### Bugs + Concerns:

Evenly distributed tests break, most likely due to loss of information from one replica. Some ideas we tried were checking tapestry nodes up to MAX_RETRIES to grab information, as well as more error checking. However, these EvenlyDist tests still fail. Any comments on potential issues/cases commonly missed and solutions would be greatly appreciated!

### Tests

To run test suite:
`go test ./test -coverpkg ./pkg/... -coverprofile=coverage.out` to test coverage on just `pkg`, then `go tool cover -html=coverage.out` to check coverage.

- `basic_io_test.go` covers basic tests to check if clients can open and close and write to files.

- `list_test.go` tests the client's `List()` function. Covers cases where list prints out nothing, list, directories, etc.

- `lock_test.go` tests the distributed locks in zookeeper

- `mkdir_test.go` tests the client's `Mkdir()` function (along with edge cases). Makes sure directories are not created under nonexistent directories, directories in general can be created and later used, etc.

- `open_test.go` tests the client's `Open()` function (along with edge cases). Makes sure opening a file follows the appropriate actions based off the arguments.

- `read_write_test.go` tests edge cases in reading and writing to files. We covered reading from a file when opening, as well as writing to it. 

- `remove_test.go` tests the client's `Remove()` function (along with edge cases). We covered removing single files, directories, and making sure they can not be accessed in the future.

- `tapestry_test.go` recreates the even distribution and load balancing tests found on Gradescope, and tests overall fault tolerance

Note: `tapestry_test.go` kills nodes that make it impossible to run with other tests, so these can be individually run! 

### Distribution of Work

Will and Mario both worked on the client interface functions and debugged them when necessary, and also both wrote numerous tests.

Will worked on the web client and setting up the gRPC server.

Mario worked on adding replication and load balancing to the tapestry nodes and data blocks.

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
