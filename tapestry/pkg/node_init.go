/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: Defines global constants and functions to create and join a new
 *  node into a Tapestry mesh, and functions for altering the routing table
 *  and backpointers of the local node that are invoked over RPC.
 */

package pkg

import (
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"google.golang.org/grpc"
)

// BASE is the base of a digit of an ID.  By default, a digit is base-16.
const BASE = 16

// DIGITS is the number of digits in an ID.  By default, an ID has 40 digits.
const DIGITS = 40

// RETRIES is the number of retries on failure. By default we have 3 retries.
const RETRIES = 3

// K is neigborset size during neighbor traversal before fetching backpointers. By default this has a value of 10.
const K = 10

// SLOTSIZE is the size each slot in the routing table should store this many nodes. By default this is 3.
const SLOTSIZE = 3

// REPUBLISH is object republish interval for nodes advertising objects.
const REPUBLISH = 10 * time.Second

// TIMEOUT is object timeout interval for nodes storing objects.
const TIMEOUT = 25 * time.Second

// Node is the main struct for the local Tapestry node. Methods can be invoked locally on this struct.
type Node struct {
	UnsafeTapestryRPCServer
	Node           RemoteNode    // The ID and address of this node
	Table          *RoutingTable // The routing table
	Backpointers   *Backpointers // Backpointers to keep track of other nodes that point to us
	LocationsByKey *LocationMap  // Stores keys for which this node is the root
	blobstore      *BlobStore    // Stores blobs on the local node
	server         *grpc.Server
}

func (local *Node) String() string {
	return fmt.Sprintf("Tapestry Node %v at %v", local.Node.ID, local.Node.Address)
}

// ID returns the tapestry node's ID in string format
func (local *Node) ID() string {
	return local.Node.ID.String()
}

// Addr returns the tapestry node's address in string format
func (local *Node) Addr() string {
	return local.Node.Address
}

// Called in tapestry initialization to create a tapestry node struct
func newTapestryNode(node RemoteNode) *Node {
	serverOptions := []grpc.ServerOption{}
	n := new(Node)

	n.Node = node
	n.Table = NewRoutingTable(node)
	n.Backpointers = NewBackpointers(node)
	n.LocationsByKey = NewLocationMap()
	n.blobstore = NewBlobStore()
	n.server = grpc.NewServer(serverOptions...)

	return n
}

// Start a node with the specified ID.
func Start(id ID, port int, connectTo string) (tapestry *Node, err error) {
	// TODO: Check if id is supplied
	// Create the RPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}

	// Get the hostname of this machine
	name, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Unable to get hostname of local machine to start Tapestry node. Reason: %v", err)
	}

	// Get the port we are bound to
	_, actualport, err := net.SplitHostPort(lis.Addr().String()) //fmt.Sprintf("%v:%v", name, port)
	if err != nil {
		return nil, err
	}

	// The actual address of this node. NOTE: If gRPC calls fail with deadline exceeded errors, this could be that it
	// is unable to resolve the computer's hostname to the local IP address. Try uncommenting the below line if this
	// happens to you (please do not check this change into your Git repo).
	// name = "127.0.0.1"
	address := fmt.Sprintf("%s:%s", name, actualport)

	// Create the local node
	tapestry = newTapestryNode(RemoteNode{ID: id, Address: address})
	fmt.Printf("Created tapestry node %v\n", tapestry)
	Trace.Printf("Created tapestry node")

	RegisterTapestryRPCServer(tapestry.server, tapestry)
	fmt.Printf("Registered RPC Server\n")
	go tapestry.server.Serve(lis)

	// If specified, connect to the provided address
	if connectTo != "" {
		// Get the node we're joining
		node, err := SayHelloRPC(connectTo, tapestry.Node)
		if err != nil {
			return nil, fmt.Errorf("Error joining existing tapestry node %v, reason: %v", address, err)
		}
		err = tapestry.Join(node)
		if err != nil {
			return nil, err
		}
	}

	return tapestry, nil
}

// Join is invoked when starting the local node, if we are connecting to an existing Tapestry.
//
// - Find the root for our node's ID
// - Call AddNode on our root to initiate the multicast and receive our initial neighbor set. Add them to our table.
// - Iteratively get backpointers from the neighbor set for all levels in range [0, SharedPrefixLength]
// - 	and populate routing table
func (local *Node) Join(otherNode RemoteNode) (err error) {
	Debug.Println("Joining", otherNode)

	// Route to our root
	root, err := local.FindRootOnRemoteNode(otherNode, local.Node.ID)
	if err != nil {
		return fmt.Errorf("Error joining existing tapestry node %v, reason: %v", otherNode, err)
	}
	// Add ourselves to our root by invoking AddNode on the remote node
	neighbors, err := root.AddNodeRPC(local.Node)
	if err != nil {
		return fmt.Errorf("Error adding ourselves to root node %v, reason: %v", root, err)
	}

	// Add the neighbors to our local routing table.
	for _, n := range neighbors {
		local.AddRoute(n)
	}

	Trace.Println("getting backpointers")

	unique := func(arr []RemoteNode) []RemoteNode {
		// TODO(TA): use NodeSet in protocol
		set := NewNodeSet()
		set.AddAll(arr)
		return set.Nodes()
	}

	// Start traversing backpointers to populate the routing table
	p := SharedPrefixLength(root.ID, local.Node.ID)
	for ; p >= 0; p-- {
		Trace.Println("backpointers level", p)
		// Trim the list to the closest unique K nodes
		neighbors = unique(neighbors)
		sort.Slice(neighbors, func(i int, j int) bool {
			return local.Node.ID.Closer(neighbors[i].ID, neighbors[j].ID)
		})
		if len(neighbors) > K {
			neighbors = neighbors[:K]
		}

		// Get the backpointers for these neighbors
		backpointers := []RemoteNode{}
		results := make(chan []RemoteNode)
		getBackpointers := func(from RemoteNode) {
			Trace.Println("Getting backpointers from", from)
			pointers, err := from.GetBackpointersRPC(local.Node, p)
			if err != nil {
				// xtr.AddTags("bad node")
				local.RemoveBadNodes([]RemoteNode{from})
			}
			results <- pointers
		}

		// Kick off the goroutines to get backpointers
		for _, node := range neighbors {
			go getBackpointers(node)
		}

		// Merge the results
		for i := 0; i < len(neighbors); i++ {
			for _, backpointer := range <-results {
				if backpointer != local.Node {
					backpointers = append(backpointers, backpointer)
				}
			}
		}

		// Add the backpointers to the routing table
		for _, b := range backpointers {
			local.AddRoute(b)
		}

		// Update the neighbors
		neighbors = append(neighbors, backpointers...)
	}

	return nil
}

// AddNode adds node to the tapestry
//
// - Begin the acknowledged multicast
// - Return the neighborset from the multicast
func (local *Node) AddNode(node RemoteNode) (neighborset []RemoteNode, err error) {
	return local.AddNodeMulticast(node, SharedPrefixLength(node.ID, local.Node.ID))
}

// AddNodeMulticast sends newNode to need-to-know nodes participating in the multicast.
// - Perform multicast to need-to-know nodes
// - Add the route for the new node (use `local.addRoute`)
// - Transfer of appropriate replica info to the new node (use `local.locationsByKey.GetTransferRegistrations`)
//   If error, rollback the location map (add back unsuccessfully transferred objects)
//
// - Propagate the multicast to the specified row in our routing table and await multicast responses
// - Return the merged neighbor set
//
// - note: `local.table.GetLevel` does not return the local node so you must manually add this to the neighbors set
func (local *Node) AddNodeMulticast(newNode RemoteNode, level int) (neighbors []RemoteNode, err error) {
	Debug.Printf("Add node multicast %v at level %v\n", newNode, level)

	// Get multicast targets
	nodes := local.Table.GetLevel(level)

	// Note: The TA implementation use recursive function calls to invoke `local.AddNodeMulticas(newNode, level + 1)`
	// 		The vanilla implementation invokes this via gRPC calls on local node as specified in the handout.
	// 		To do that, local node must be appended to list of nodes that are notified.
	// 		And the recursive call below should be removed.
	//  	However, vanilla implementation is slower (~2x) and vulnerable to network errors. Students are encouraged
	//		to think about non-vanilla implementation.
	//
	// Note: Uncomment for vanilla implementation of `AddNodeMulticast`
	// nodes = append(nodes, local.node)

	done := make(chan []RemoteNode)

	notify := func(destination RemoteNode) {
		Trace.Println("Notifying", destination)
		newNeighbors, err := destination.AddNodeMulticastRPC(newNode, level+1)
		if err != nil {
			local.RemoveBadNodes([]RemoteNode{destination})
		}
		done <- newNeighbors
	}

	// Kick off asynchronous multicast
	for _, node := range nodes {
		go notify(node)
	}

	// If we're at level DIGITS, transfer keys, otherwise multicast to 1 level down
	if level == DIGITS {
		// Transfer keys to the new node
		go func() {
			// xtr.NewTask("transferkeys")
			Trace.Print("Beginning transfer keys")
			// Add the new node to the routing table
			local.AddRoute(newNode)

			// Get the data to transfer
			objects := local.LocationsByKey.GetTransferRegistrations(local.Node, newNode)

			// Transfer the data
			err := newNode.TransferRPC(local.Node, objects)

			if err != nil {
				// On error, remove the new node from our table, and reinsert the transferred data
				local.RemoveBadNodes([]RemoteNode{newNode})
				local.LocationsByKey.RegisterAll(objects, TIMEOUT)
			}
		}()

		neighbors = append(neighbors, local.Node)
	} else {
		// Multicast the local node at the next level down

		// Note: comment this recursive call for vanilla implementation of `AddNodeMulticast`
		neighbors, err = local.AddNodeMulticast(newNode, level+1)

		// Append returned neighbors
		for i := 0; i < len(nodes); i++ {
			neighbors = append(neighbors, <-done...)
		}
	}

	return
}

// AddBackpointer adds the from node to our backpointers, and possibly add the node to our
// routing table, if appropriate
func (local *Node) AddBackpointer(from RemoteNode) (err error) {
	if local.Backpointers.Add(from) {
		Debug.Printf("Added backpointer %v\n", from)
	}
	local.AddRoute(from)
	return
}

// RemoveBackpointer removes the from node from our backpointers
func (local *Node) RemoveBackpointer(from RemoteNode) (err error) {
	if local.Backpointers.Remove(from) {
		Debug.Printf("Removed backpointer %v\n", from)
	}
	return
}

// GetBackpointers gets all backpointers at the level specified, and possibly add the node to our
// routing table, if appropriate
func (local *Node) GetBackpointers(from RemoteNode, level int) (backpointers []RemoteNode, err error) {
	Debug.Printf("Sending level %v backpointers to %v\n", level, from)
	backpointers = local.Backpointers.Get(level)
	local.AddRoute(from)
	return
}

// RemoveBadNodes discards all the provided nodes
// - Remove each node from our routing table
// - Remove each node from our set of backpointers
func (local *Node) RemoveBadNodes(badnodes []RemoteNode) (err error) {
	for _, badnode := range badnodes {
		if local.Table.Remove(badnode) {
			Debug.Printf("Removed bad node %v\n", badnode)
		}
		if local.Backpointers.Remove(badnode) {
			Debug.Printf("Removed bad node backpointer %v\n", badnode)
		}
	}
	return
}

// Utility function that adds a node to our routing table.
//
// - Adds the provided node to the routing table, if appropriate.
// - If the node was added to the routing table, notify the node of a backpointer
// - If an old node was removed from the routing table, notify the old node of a removed backpointer
func (local *Node) AddRoute(node RemoteNode) (err error) {
	added, previous := local.Table.Add(node)
	go func() {
		if previous != nil {
			// Notify of backpointer removal. Error doesn't matter because previous is no longer in table
			previous.RemoveBackpointerRPC(local.Node)
		}

		if added {
			Debug.Printf("Added %v to routing table\n", node)
			// Try notifying of a backpointer
			err := node.AddBackpointerRPC(local.Node)

			if err != nil {
				Debug.Printf("Backpointer notification to %v failed\n", node)
				// If backpointer notification fails, remove the node from the routing table
				// Note: this line is necessary to pass TopologySuite.TestTA_AddRoute5_Eviction_Notification_Failure_Test
				local.Table.Remove(node)

				// Attempt reinsertion of previous
				if previous != nil {
					local.AddRoute(*previous)
				}
			}
		}
	}()

	return
}
