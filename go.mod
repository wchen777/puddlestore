module puddlestore

go 1.17

require (
	github.com/go-zookeeper/zk v1.0.2
	github.com/hashicorp/go-msgpack v1.1.5
	tapestry v0.0.0-00010101000000-000000000000
)

require github.com/google/uuid v1.3.0

// Replace this
replace tapestry => ./tapestry // /path/to/your/tapestry/implementation/root/folder
