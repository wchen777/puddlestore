package pkg

import (
	"github.com/go-zookeeper/zk"
)

// Client is a puddlestore client interface that will communicate with puddlestore nodes
type Client interface {
	// `Open` opens a file and returns a file descriptor. If the `create` is true and the
	// file does not exist, create the file. If `create` is false and the file does not exist,
	// return an error. If `write` is true, then flush the resulting inode on Close(). If `write`
	// is false, no need to flush the inode to zookeeper. If `Open` is successful, the returned
	// file descriptor should be unique to the file. The client is responsible for keeping
	// track of local file descriptors. Using `Open` allows for file-locking and
	// multi-operation transactions.
	Open(path string, create, write bool) (int, error)

	// `Close` closes the file and flushes its contents to the distributed filesystem.
	// The updated closed file should be able to be opened again after successfully closing it.
	// We only flush changes to the file on close to ensure copy-on-write atomicity of operations.
	// Refer to the handout for more information on why this is necessary.
	Close(fd int) error

	// `Read` returns a `size` amount of bytes starting at `offset` in an opened file.
	// Reading at non-existent offset returns empty buffer and no error.
	// If offset+size exceeds file boundary, return as much as possible with no error.
	// Returns err if fd is not opened.
	Read(fd int, offset, size uint64) ([]byte, error)

	// `Write` writes `data` starting at `offset` on an opened file. Writing beyond the
	// file boundary automatically fills the file with zero bytes. Returns err if fd is not opened.
	// If the file was opened with write = true flag, `Write` should return an error.
	Write(fd int, offset uint64, data []byte) error

	// `Mkdir` creates directory at the specified path.
	// Returns error if any parent directory does not exist (non-recursive).
	Mkdir(path string) error

	// `Remove` removes a directory or file. Returns err if not exists.
	Remove(path string) error

	// `List` lists file & directory names (not full names) under `path`. Returns err if not exists.
	List(path string) ([]string, error)

	// Release zk connection. Subsequent calls on Exit()-ed clients should return error.
	Exit()
}

// TODO: implement the Client interface
type PuddleClient struct {
	ID     string
	zkConn *zk.Conn
}

// open a file and return a file descriptor
func (c *PuddleClient) Open(path string, create, write bool) (int, error) {
	return 0, nil
}

// close a file and flush its contents to the distributed filesystem
func (c *PuddleClient) Close(fd int) error {
	return nil
}

// read a file and return a buffer of size `size` starting at `offset`
func (c *PuddleClient) Read(fd int, offset, size uint64) ([]byte, error) {
	return nil, nil
}

// write a file and write `data` starting at `offset`
func (c *PuddleClient) Write(fd int, offset uint64, data []byte) error {
	return nil
}

// create a directory at the specified path
func (c *PuddleClient) Mkdir(path string) error {
	return nil
}

// remove a directory or file
func (c *PuddleClient) Remove(path string) error {
	return nil
}

// list file & directory names (not full names) under `path`
func (c *PuddleClient) List(path string) ([]string, error) {
	return nil, nil
}

// release zk connection
func (c *PuddleClient) Exit() {
	c.zkConn.Close()
}
