package pkg

import (
	// "fmt"
	// "strings"

	// "sort"

	"github.com/go-zookeeper/zk"
)

// IDEAS FOR LOCKS:

/*
 - every inode has a corresponding dist lock that is a znode in the locks subtree (e.g. /locks/inode-1)
 - the dist lock uses the file path (without the fs root prefix and instead the locks prefix) as the znode name
 - each lock has a read lock and write lock for granularity
	- how do we want to accomplish this?
	- use diff znodes for each lock?
	- we would need to use waits, etc. to notify those who hold read locks that a write lock has been acquired
	- how do we grab a write lock when others hold read locks
*/

// FROM ED:

/*
Each file has one write-lock and infinite read-locks.
To write you must get a write-lock, while a client has a write lock, no other client can hold a read-lock.

When client(s) have a read-lock, to acquire a write-lock,
the system must either wait until all read-locks are returned or it may force the clients
to give up their read locks (a system designed at Microsoft did this).
*/

/*
Implementation designs:

Place locks on Zookeeper paths that contain inodes.
In your lock implementation, be careful about filtering out lock-unrelated znodes;
you might not have done this in the Zookeeper lab. TODO: WHAT DOES THIS IMPLY??? is this about directories?

Put all the locks together on their own root in Zookeeper, e.g. all under /locks.
This avoids having to filter out znodes that store inodes. WE ALREADY DO THIS.

When a client with a lock crashes, it will be unable to release the lock,
the system must be able to detect the failure and reclaim the locks. TODO: this.
*/

type DistLock struct {
	zkConn   *zk.Conn
	lockRoot string // e.g. /locks
	path     string // full zk path of the lock

}

// CreateDistLock creates a distributed lock
func CreateDistLock(root string, zkConn *zk.Conn) *DistLock {
	dlock := &DistLock{
		lockRoot: root,
		path:     "",
		zkConn:   zkConn,
	}
	return dlock
}

// grab a read lock
func (d *DistLock) RLock() (err error) {
	return nil
}

// grab a write lock
func (d *DistLock) WLock() (err error) {
	return nil
}

// unlock read lock
func (d *DistLock) RUnlock() (err error) {
	return nil
}

// unlock write lock
func (d *DistLock) WUnlock() (err error) {
	return nil
}
