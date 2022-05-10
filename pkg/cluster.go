package pkg

import (
	"errors"
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
const SEED = 4444

// TODO: HASHING FUNCTION FOR LOAD BALANCING, consistent hashing, round robin, etc.
func (c *Cluster) HashClientIDToTapNode(clientID string) *Tapestry {

	return c.nodes[0] // PLACEHOLDER
}

// NewClient creates a new Puddlestore client
func (c *Cluster) NewClient() (Client, error) {
	// TODO: Return a new PuddleStore Client that implements the Client interface
	// todo pt2: we should be listing out ALL tapestry clients
	// tapestry/node-01-addr
	// tapestry/node-02-addr
	// when a client connects, we select from these tap clients with load balencing
	// each client isn't strictly assigned a tap node, we simply just use these tap clients for accessing data
	// see ed post #649 + fault tolerence section in gearup.

	// try and establish a new connection
	conn, err := ConnectZk(c.config.ZkAddr)

	if err != nil {
		return nil, err
	}

	clientID := uuid.New()

	client := &PuddleClient{
		ID:     clientID.String(),
		zkConn: conn,
		fsPath: FS_ROOT,
		// locksPath:    LOCK_ROOT,
		tapestryPath: TAP_ADDRESS_ROOT,

		openFiles:  make([]*OpenFile, CLIENT_OPEN_FILES_LIMIT),
		dirtyFiles: make(map[int]bool),
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateCluster starts all nodes necessary for puddlestore
func CreateCluster(config Config) (*Cluster, error) {

	// try and establish a new connection
	conn, err := ConnectZk(config.ZkAddr)

	if err != nil {
		return nil, err
	}

	// create tapestry root path
	conn.Create(TAP_ADDRESS_ROOT, []byte{}, 0, zk.WorldACL(zk.PermAll))

	// init paths directories: /puddlestore, /locks, /tapestry
	initPaths(conn)

	// create random set of tapestries of count config.NumTapestry
	randNodes, err := tapestry.MakeRandomTapestries(SEED, config.NumTapestry)

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

	// register nodes in /tapestry ...
	for _, node := range nodes {

		// encode tap node
		Tapinode := &TapestryAddrNode{
			Addr: node.tap.Node.Address, // contains tap address to connect to
		}

		// encode a inode with tap node address.
		tapNodeBuffer, err := encodeMsgPack(Tapinode)

		if err != nil {
			return nil, err
		}

		// what if tapestry nodes fail?
		_, err = conn.Create(TAP_ADDRESS_ROOT+"/node-", tapNodeBuffer.Bytes(), zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

		if err != nil {
			return nil, errors.New("failed to create tapestry node: " + err.Error())
		}

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

// initializes the zookeeper internal file system and locks directory paths
func initPaths(c *zk.Conn) error {

	// create the inode
	newFileinode := &inode{
		Filepath: FS_ROOT,              // this is the path of the file in the actual filesystem
		Size:     0,                    // this is the size of the file in bytes (starts as empty)
		Blocks:   make([]uuid.UUID, 0), // this is the list of data blocks (each block is a uuid that represents an entry in tapestry)
		IsDir:    true,                 // this is the flag that indicates if the file is a directory
	}

	// marshal the inode to bytes
	inodeBuffer, err := encodeInode(*newFileinode)

	if err != nil {
		return err
	}

	// if fs path does not exist, create it
	_, err = c.Create(FS_ROOT, inodeBuffer, 0, zk.WorldACL(zk.PermAll))

	if err != nil {
		return err
	}

	_, err = c.Create(TAP_ADDRESS_ROOT, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	_, err = c.Create(LOCK_ROOT, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	return nil
}

// RANDOM IDEAS:

// IMPLEMENT LOAD BALANCING SO ZKCONN IS ASSIGNED TO A DIFF TAP NODE?

// IDEAS:
// - generate new uuid for client
// - use consistent hashing to assign client to tapestry node using uuid
