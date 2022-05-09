package pkg

import (
	tapestry "tapestry/pkg"

	"github.com/go-zookeeper/zk"
)

// Tapestry is a wrapper for a single Tapestry node. It is responsible for
// maintaining a zookeeper connection and implementing methods we provide
type Tapestry struct {
	tap *tapestry.Node
	zk  *zk.Conn
}

type TapestryAddrNode struct {
	Addr string
}

// NewTapestry creates a new tapestry struct
func NewTapestry(tap *tapestry.Node, zkAddr string) (*Tapestry, error) {
	//  create new zkConn.
	zkConn, err := ConnectZk(zkAddr)

	if err != nil {
		return nil, err
	}

	// returns back tapestry node with zkConn in struct.
	return &Tapestry{tap: tap, zk: zkConn}, nil
}

// GracefulExit closes the zookeeper connection and gracefully shuts down the tapestry node
func (t *Tapestry) GracefulExit() {
	t.zk.Close()
	t.tap.Leave()
}
