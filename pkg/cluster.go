package pkg

import (
	tapestry "tapestry/pkg"
	"time"

	"github.com/go-zookeeper/zk"
	uuid "github.com/google/uuid"
)

// Cluster is an interface for all nodes in a puddlestore cluster. One should be able to shutdown
// this cluster and create a client for this cluster
type Cluster struct {
	config Config
	nodes  []*Tapestry
}

const LOCK_ROOT = "/locks"
const FS_ROOT = "/puddlestore"
const TAP_ADDRESS_ROOT = "/tapestry"

const CLIENT_OPEN_FILES_LIMIT = 256

// TODO: HASHING FUNCTION FOR LOAD BALANCING, consistent hashing, round robin, etc.
func (c *Cluster) HashClientIDToTapNode(clientID string) *Tapestry {

	return c.nodes[0] // PLACEHOLDER
}

// NewClient creates a new Puddlestore client
func (c *Cluster) NewClient() (Client, error) {
	// TODO: Return a new PuddleStore Client that implements the Client interface

	// try and establish a new connection
	conn, err := ConnectZk(c.config.ZkAddr)

	if err != nil {
		return nil, err
	}

	clientID := uuid.New()

	client := &PuddleClient{
		ID:           clientID.String(),
		zkConn:       conn,
		fsPath:       FS_ROOT,
		locksPath:    LOCK_ROOT,
		tapestryPath: TAP_ADDRESS_ROOT,

		openFiles: make([]*OpenFile, 256),
	}

	// init paths for all zk roots, if they do not exist
	err = client.initPaths()

	if err != nil {
		return nil, err
	}

	// TODO: WHAT SHOULD WE STORE IN ZK, JUST THE TAP ADDR OR THE WHOLE TAP NODE OBJECT?
	assignedTapNode := c.HashClientIDToTapNode(client.ID)    // use load balancing to assign client to tapestry node
	tapNodeMarshalled, err := encodeMsgPack(assignedTapNode) // marshal tapestry node to be stored in zk

	if err != nil {
		return nil, err
	}

	// use the zk conn to create a new sequential + ephemeral node to represent the client
	// the data stored at this znode to hold the assigned tapestry address (RN IT IS WHOLE TAP NODE OBJ)
	_, err = client.zkConn.Create(client.tapestryPath+"/"+client.ID, tapNodeMarshalled.Bytes(), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

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
