package pkg

import (
	tapestry "tapestry/pkg"
	"time"
)

// Cluster is an interface for all nodes in a puddlestore cluster. One should be able to shutdown
// this cluster and create a client for this cluster
type Cluster struct {
	config Config
	nodes  []*Tapestry
}

const LOCK_ROOT = "/locks"
const FS_ROOT = "/puddlestore"

// NewClient creates a new Puddlestore client
func (c *Cluster) NewClient() (Client, error) {
	// TODO: Return a new PuddleStore Client that implements the Client interface

	// try and establish a new connection
	conn, err := ConnectZk(c.config.ZkAddr)

	if err != nil {
		return nil, err
	}

	client := &PuddleClient{
		ID:        "", // TODO: what should we do for ID?
		zkConn:    conn,
		fsPath:    FS_ROOT,
		locksPath: LOCK_ROOT,
	}

	err = client.initPaths()

	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateCluster starts all nodes necessary for puddlestore
func CreateCluster(config Config) (*Cluster, error) {

	// TODO: Start your tapestry cluster with size config.NumTapestry. You should
	// also use the zkAddr (zookeeper address) found in the config and pass it to
	// your Tapestry constructor method

	// create random set of tapestries of count config.NumTapestry
	randNodes, err := tapestry.MakeRandomTapestries(4444, config.NumTapestry)

	var nodes []*Tapestry

	// iterate through newly created nodes to create Tapestry Nodes for *Cluster
	for i := 0; i < config.NumTapestry; i += 1 {

		nodeToAdd, err := NewTapestry(randNodes[i], config.ZkAddr)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, nodeToAdd)
	}

	if err != nil {
		return nil, err
	}

	return &Cluster{config: config, nodes: nodes}, nil
}

// Shutdown causes all tapestry nodes to gracefully exit
func (c *Cluster) Shutdown() {
	for _, node := range c.nodes {
		node.GracefulExit()
	}

	time.Sleep(time.Second)
}

// RANDOM IDEAS:

// IMPLEMENT LOAD BALANCING SO ZKCONN IS ASSIGNED TO A DIFF TAP NODE?

// IDEAS:
// - generate new uuid for client
// - use consistent hashing to assign client to tapestry node using uuid
