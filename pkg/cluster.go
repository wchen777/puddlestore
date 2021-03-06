package pkg

import (
	"errors"
	tapestry "tapestry/pkg"
	"time"

	"github.com/go-zookeeper/zk"
	uuid "github.com/google/uuid"
)

// Client is a puddlestore client interface that will communicate with puddlestore nodes
var MAX_RETRIES = 3

// Cluster is an interface for all nodes in a puddlestore cluster. One should be able to shutdown
// this cluster and create a client for this cluster
type Cluster struct {
	config Config
	nodes  []*Tapestry
}

const LOCK_ROOT = "/locks"
const FS_ROOT = "/puddlestore"
const TAP_ADDRESS_ROOT = "/tapestry"

const CLIENT_OPEN_FILES_LIMIT = 2048
const SEED = 4444

var BLOCK_SIZE uint64 // will be overwritten by config file.

// NewClient creates a new Puddlestore client
func (c *Cluster) NewClient() (Client, error) {
	// each client isn't strictly assigned a tap node, we simply just use these tap clients for accessing data
	// see ed post #649 + fault tolerence section in gearup.

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
		numReplicas:  c.config.NumReplicas,
		tapestryPath: TAP_ADDRESS_ROOT,

		openFiles:  make([]*OpenFile, CLIENT_OPEN_FILES_LIMIT),
		dirtyFiles: make(map[int]bool),
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

// FOR TESTING, ADJUST PARAMETERS AS NECESSARY
func (c *Cluster) NewClientTest(maxFiles int) (Client, error) {
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
		numReplicas:  c.config.NumReplicas,
		tapestryPath: TAP_ADDRESS_ROOT,

		openFiles:  make([]*OpenFile, maxFiles), // adjust max open files to be small
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

	BLOCK_SIZE = config.BlockSize

	if err != nil {
		return nil, err
	}

	// create tapestry root path
	//conn.Create(TAP_ADDRESS_ROOT, []byte{}, 0, zk.WorldACL(zk.PermAll))

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
			Id:     node.tap.Node.ID.String(),
			Addr:   node.tap.Node.Address, // contains tap address to connect to
			TapCli: nil,
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

	// try and establish a new connection
	conn, err := ConnectZk(c.config.ZkAddr)

	if err != nil {
		return
	}

	// delete puddlestore
	children, _, _ := conn.Children(FS_ROOT)
	for _, c := range children {
		conn.Delete(FS_ROOT+"/"+c, -1)
	}
	conn.Delete(FS_ROOT, -1)

	// delete tap
	children, _, _ = conn.Children(TAP_ADDRESS_ROOT)
	for _, c := range children {
		conn.Delete(TAP_ADDRESS_ROOT+"/"+c, -1)
	}
	conn.Delete(TAP_ADDRESS_ROOT, -1)

	// delete locks
	children, _, _ = conn.Children(LOCK_ROOT)
	for _, c := range children {
		conn.Delete(LOCK_ROOT+"/"+c, -1)
	}
	conn.Delete(LOCK_ROOT, -1)

	conn.Close()
	time.Sleep(time.Second)
}

// return the list of tap nodes (for testing downing tap nodes)
func (c *Cluster) GetTapestryNodes() []*Tapestry {
	return c.nodes
}

// initializes the zookeeper internal file system and locks directory paths
func initPaths(c *zk.Conn) error {

	// create the inode
	newFileinode := &inode{
		Filepath: FS_ROOT,           // this is the path of the file in the actual filesystem
		Size:     0,                 // this is the size of the file in bytes (starts as empty)
		Blocks:   make([]string, 0), // this is the list of data blocks (each block is a uuid that represents an entry in tapestry)
		IsDir:    true,              // this is the flag that indicates if the file is a directory
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
