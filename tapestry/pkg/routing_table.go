/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: Defines the RoutingTable type and provides methods for interacting
 *  with it.
 */

package pkg

import (
	"sync"
)

// RoutingTable has a number of levels equal to the number of digits in an ID
// (default 40). Each level has a number of slots equal to the digit base
// (default 16). A node that exists on level n thereby shares a prefix of length
// n with the local node. Access to the routing table protected by a mutex.
type RoutingTable struct {
	local RemoteNode                 // The local tapestry node
	Rows  [DIGITS][BASE][]RemoteNode // The rows of the routing table
	mutex sync.Mutex                 // To manage concurrent access to the routing table (could also have a per-level mutex)
}

// NewRoutingTable creates and returns a new routing table, placing the local node at the
// appropriate slot in each level of the table.
func NewRoutingTable(me RemoteNode) *RoutingTable {
	t := new(RoutingTable)
	t.local = me

	// Create the node lists with capacity of SLOTSIZE
	for i := 0; i < DIGITS; i++ {
		for j := 0; j < BASE; j++ {
			t.Rows[i][j] = make([]RemoteNode, 0, SLOTSIZE)
		}
	}

	// Make sure each row has at least our node in it
	for i := 0; i < DIGITS; i++ {
		slot := t.Rows[i][t.local.ID[i]]
		t.Rows[i][t.local.ID[i]] = append(slot, t.local)
	}

	return t
}

// Add adds the given node to the routing table.
//
// Note you should not add the node to preceding levels. You need to add the node
// to one specific slot in the routing table (or replace an element if the slot is full
// at SLOTSIZE).
//
// Returns true if the node did not previously exist in the table and was subsequently added.
// Returns the previous node in the table, if one was overwritten.
func (t *RoutingTable) Add(node RemoteNode) (added bool, previous *RemoteNode) {
	// Check we aren't re-adding ourselves
	if t.local.ID == node.ID {
		return
	}

	// Get the level of the table where this node should go
	level := t.level(node)
	digit := node.ID[level]

	t.mutex.Lock()
	defer t.mutex.Unlock()

	slot := t.Rows[level][digit]
	if len(slot) == SLOTSIZE {
		added, previous = doReplace(t.local, node, slot)
	} else {
		slot, added = doAdd(node, slot)
	}
	t.Rows[level][digit] = slot

	return
}

// Remove removes the specified node from the routing table, if it exists.
// Returns true if the node was in the table and was successfully removed.
// Return false if a node tries to remove itself from the table.
func (t *RoutingTable) Remove(node RemoteNode) (wasRemoved bool) {
	// Cannot remove ourselves from the table
	if t.local == node {
		return false
	}

	// Determine the level and slot the node belongs in
	level := t.level(node)
	digit := node.ID[level]

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.Rows[level][digit], wasRemoved = doRemove(node, t.Rows[level][digit])

	return
}

// GetLevel gets ALL nodes on the specified level of the routing table, EXCLUDING the local node.
func (t *RoutingTable) GetLevel(level int) (nodes []RemoteNode) {
	if level < 0 || level >= DIGITS {
		return nil
	}

	nodes = make([]RemoteNode, 0, BASE*SLOTSIZE)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, slot := range t.Rows[level] {
		for _, node := range slot {
			if node != t.local {
				nodes = append(nodes, node)
			}
		}
	}

	return
}

// FindNextHop searches the table for the closest next-hop node for the provided ID starting at the given level.
func (t *RoutingTable) FindNextHop(id ID, level int32) RemoteNode {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for ; level < DIGITS-1; level++ {
		node := t.doGetNodeAtLevel(level, id)
		if node != t.local {
			return node
		}
	}

	return t.local
}

// Private non-locking implementation.
func (t *RoutingTable) doGetNodeAtLevel(d int32, id ID) (node RemoteNode) {
	// Get the d'th row, then cycle through slots until we find a node
	row := t.Rows[d]
	digit := id[d]
	for i := 0; i < BASE; i++ {
		slot := row[digit]
		if len(slot) > 0 {
			return closest(id, slot)
		}
		digit = (digit + 1) % BASE
	}

	return t.local
}

func (t *RoutingTable) level(node RemoteNode) int {
	return SharedPrefixLength(t.local.ID, node.ID)
}

// Removes all occurrences of toRemove from nodes.
func doRemove(toRemove RemoteNode, nodes []RemoteNode) ([]RemoteNode, bool) {
	wasRemoved := false
	size := len(nodes)
	for i := 0; i < size; i++ {
		if nodes[i] == toRemove {
			lastnode := nodes[size-1]
			nodes[size-1] = toRemove
			nodes[i] = lastnode
			nodes = nodes[:size-1]
			i--
			wasRemoved = true
			size--
		}
	}
	return nodes, wasRemoved
}

// If the new node is closer than an existing node, the existing node is replaced.
func doReplace(local RemoteNode, newNode RemoteNode, existingNodes []RemoteNode) (existingNodeWasReplaced bool, previous *RemoteNode) {
	// First, check the node isn't already in the list
	for i := 0; i < len(existingNodes); i++ {
		if (existingNodes)[i] == newNode {
			return false, nil
		}
	}

	// Now, try replacing an existing node with the new node
	furthest := newNode
	for i := 0; i < len(existingNodes); i++ {
		existing := existingNodes[i]
		if local.ID.Closer(furthest.ID, existing.ID) {
			existingNodes[i] = furthest
			furthest = existing
			existingNodeWasReplaced = true
		}
	}
	if furthest != newNode {
		previous = &furthest
	}
	return
}

// Add a node to the list so long as it's not already present
func doAdd(newNode RemoteNode, existingNodes []RemoteNode) ([]RemoteNode, bool) {
	for i := 0; i < len(existingNodes); i++ {
		if existingNodes[i] == newNode {
			return existingNodes, false
		}
	}
	existingNodes = append(existingNodes, newNode)
	return existingNodes, true
}

// Returns the closest node in the list to the provided ID
func closest(id ID, nodes []RemoteNode) (closest RemoteNode) {
	closest = nodes[0]
	for _, node := range nodes {
		if id.Closer(node.ID, closest.ID) {
			closest = node
		}
	}
	return
}
