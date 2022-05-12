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

### Some expectations for using this README

Document your tests, known bugs, and extra features here.

### Tests

To run test suite:
`go test ./test -coverpkg ./pkg/... -coverprofile=coverage.out` to test coverage on just `pkg`, then `go tool cover -html=coverage.out` to check coverage.


### How to Run Web Client

1. Switch to `web-client` branch
2. Ensure Zookeeper is running on port 2181
3. Run the server using `go run api/main.go` from the root.
4. Run the following command: `./grpcwebproxy-v0.15.0-osx-x86_64 --backend_addr=:3333 --server_http_debug_port 3334 --allow_all_origins --run_tls_server=false` in the root to ensure the proxy forwards connections from the client to the server.
5. Run the client either by running `npm run start` from within the `puddlestore-web-client` directory, or opening the `index.html` in the `puddlestore-web-client-build` directory (TODO).
6. Have fun!

![puddlestore client](img/puddlestore%20sc.png "PuddleStore Client")

`protoc puddlestore.proto \
--js_out=import_style=commonjs,binary:./src/puddlestore \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/puddlestore`
