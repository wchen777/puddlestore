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

`./grpcwebproxy-v0.15.0-osx-x86_64 --backend_addr=:3333 --server_http_debug_port 3334 --allow_all_origins --run_tls_server=false`

`protoc puddlestore.proto \
--js_out=import_style=commonjs,binary:./src/puddlestore \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/puddlestore
`
