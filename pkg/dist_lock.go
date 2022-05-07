package pkg

import (
	// "fmt"
	// "strings"

	// "sort"

	"fmt"
	"sort"
	"strings"

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

const lockPrefix = "/locks-"

type DistLock struct {
	zkConn   *zk.Conn
	lockRoot string // e.g. /a
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

// NORMAL LOCK IMPLEMENTATION:
// todo: realease lock if crash
func (d *DistLock) Acquire() (err error) {

	// 1. call create
	// locks /dir/lock-
	path, err := d.zkConn.Create(d.lockRoot+lockPrefix, []byte(""), zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	fmt.Println("path: ", path)

	if err != nil {
		return err
	}

	// update path field of distlock
	d.path = path

	for {

		// call children on lock node without watch flag
		// puddlestore concern:
		// children on \a with dir \a\b will give us:
		// \a\lock-..., as well as \a\b elements...
		unfiltChil, _, err := d.zkConn.Children(d.lockRoot)

		// filter non lock children
		var chil []string

		for _, s := range unfiltChil {

			// if lock root AND has lock prefix, add to new list
			// ignore children that are NOT locks of root dir(root + \lock-).
			if strings.HasPrefix(s, (d.lockRoot + lockPrefix)) {
				chil = append(chil, s)
			}
		}

		if err != nil {
			return err
		}

		// sort children array by string comparison
		sort.Strings(chil)

		// if pathname in step 1 has lowest squence number suffix, client has lock and should exit protocol
		if len(chil) <= 1 || strings.HasSuffix(d.path, chil[0]) {
			// client has lock and should exit protocol
			return nil
		}

		// call exists with watch flag on path in lock directory with next lowest sequence number
		// get index of path suffix in children array

		index := sort.SearchStrings(chil, path[len(d.lockRoot)+1:])

		exists, _, eventChan, err := d.zkConn.ExistsW(d.lockRoot + "/" + chil[index-1])

		if err != nil {
			return err
		}

		if exists {
			// wait for notification for pathname from previous step

			<-eventChan
		}

	}

}

// The unlock protocol is very simple: clients wishing to release a lock simply delete the node they created in step 1.
func (d *DistLock) Release() (err error) {

	return d.zkConn.Delete(d.path, -1)
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
