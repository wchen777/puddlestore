package pkg

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	tapestry "tapestry/pkg"

	"github.com/go-zookeeper/zk"
	uuid "github.com/google/uuid"
)

// Client is a puddlestore client interface that will communicate with puddlestore nodes

// each block is 4000 bytes.(4kb)
const BLOCK_SIZE = 4096

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

type OpenFile struct {
	INode    *inode
	Data     []byte    // cached file block data, TODO: how do we get all the data in one place from tapestry blocks ?
	FileLock *DistLock // file lock
}

type PuddleClient struct {
	ID        string
	zkConn    *zk.Conn
	openFiles []*OpenFile // map from file descriptor to inode, represented as an array of inodes (each fd is an index in the array)

	// lock states field here, best way to store read and write locks associated with an inode?

	fsPath       string // file system path prefix within zookeeper, e.g. /puddlestore
	tapestryPath string // path root for tapestry nodes assigned to each client, e.g. /tapestry

	dirtyFiles map[int]bool // dirty files set ? for flushing, should contain file descriptors (or file paths)?

}

// ---------------------- CLIENT INTERFACE IMPLEMENTATION ---------------------- //

// open a file and return a file descriptor, DOES THIS PATH START WITH A /?
func (c *PuddleClient) Open(path string, create, write bool) (int, error) {

	// search for the file path metadata in zookeeper
	fileExists, _, err := c.zkConn.Exists(c.fsPath + "/" + path)

	if err != nil {
		return -1, err
	}

	// create lock for file
	distlock := CreateDistLock(c.fsPath+"/"+path, c.zkConn)

	distlock.Acquire()

	var newFileinode *inode
	data := make([]byte, 0)

	if !fileExists { // if the file metadata does not exist in the zookeeper fs

		if !create { // if we are not creating and the file does not exist, return error
			distlock.Release()
			return -1, zk.ErrNoNode
		} else { // otherwise create the file

			// create the inode
			newFileinode = &inode{
				Filepath: path,                 // this is the path of the file in the actual filesystem
				Size:     0,                    // this is the size of the file in bytes (starts as empty)
				Blocks:   make([]uuid.UUID, 0), // this is the list of data blocks (each block is a uuid that represents an entry in tapestry)
			}

			// marshal the inode to bytes
			inodeBuffer, err := encodeInode(*newFileinode)

			if err != nil { // encode fails
				// release the lock
				distlock.Release()
				return -1, err
			}

			// create the file metadata in zookeeper, should be neither sequential nor ephemeral
			c.zkConn.Create(c.fsPath+"/"+path, inodeBuffer, 0, zk.WorldACL(zk.PermAll))

		}

	} else {

		// get the inode from zookeeper
		data, _, err := c.zkConn.Get(c.fsPath + "/" + path)

		if err != nil {
			distlock.Release()
			return -1, err
		}

		// unmarshal the inode
		newFileinode, err = decodeInode(data)
		if err != nil {
			distlock.Release()
			return -1, err
		}

		// TODO: READ THE FILE DATA FROM TAPESTRY USING BLOCKS FOUND IN INODE
		data = make([]byte, 0)

	}

	// get next client file descriptor
	fd := c.findNextFreeFD()

	if fd == -1 { // if there are no free file descriptors, return error
		distlock.Release()
		return -1, errors.New("no free file descriptors, ENOMEM")
	}

	// add the file to the open files list
	c.openFiles[fd] = &OpenFile{
		INode:    newFileinode,
		Data:     data,
		FileLock: distlock,
	}

	// if we have specified write, add fd to dirty files (to be flushed on close)
	if write {
		c.dirtyFiles[fd] = true
	}

	return fd, nil

}

// close a file and flush its contents to the distributed filesystem
func (c *PuddleClient) Close(fd int) error {

	// open file
	openFile := c.openFiles[fd]

	// check dirty files set
	if c.dirtyFiles[fd] {
		// flush the file
		selectedNode, err := c.getRandomTapestryNode() // TODO: check this logic in this helper

		if err != nil {
			// release lock.
			openFile.FileLock.Release()

			// close conn
			c.zkConn.Close()

			return err
		}

		toDecode, _, err := c.zkConn.Get(selectedNode)

		if err != nil {
			// release lock.
			openFile.FileLock.Release()

			// close conn
			c.zkConn.Close()

			return err
		}

		inode, err := decodeInode(toDecode)

		if err != nil {
			// release lock.
			openFile.FileLock.Release()

			// close conn
			c.zkConn.Close()

			return err
		}

		// connects to tap belonging to inode.Filepath(which is addr of tap node)
		client, err := tapestry.Connect(inode.Filepath)

		if err != nil {
			// release lock.
			openFile.FileLock.Release()

			// close conn
			c.zkConn.Close()

			return err
		}

		// release lock.
		openFile.FileLock.Release()

		// keeps track of end of array to get correct slice of bytes.
		var end int

		// keeps track of new uuids
		var newUIDs []uuid.UUID

		for i := 0; i < len(openFile.Data); i += BLOCK_SIZE {

			end += BLOCK_SIZE

			// prevents slice beyond boundary.
			if end > len(openFile.Data) {
				end = len(openFile.Data)
			}

			// new 4kb data
			dataBlock := openFile.Data[i:end]

			// create new uuid, store into tapestry uuid associated with block
			newUID := uuid.New()

			client.Store(newUID.String(), dataBlock)

			// add to array of newuids to replace old in inode.
			newUIDs = append(newUIDs, newUID)
		}

		// here: done populating tap with newuids <--> blocks
		// replace inode uids with new
		openFile.INode.Blocks = newUIDs

		// clear fd
		c.openFiles[fd] = nil

		// remove the file from the dirty files set
		delete(c.dirtyFiles, fd)
	}

	return nil
}

// read a file and return a buffer of size `size` starting at `offset`
func (c *PuddleClient) Read(fd int, offset, size uint64) ([]byte, error) {

	// TODO: for offsets, create a helper function to get the correct data blocks from offset
	return nil, nil
}

// write a file and write `data` starting at `offset`
func (c *PuddleClient) Write(fd int, offset uint64, data []byte) error {

	// remember to modify the inode data stored locally on each write, flush to zookeeper on close
	return nil
}

// create a directory at the specified path
func (c *PuddleClient) Mkdir(path string) error {
	return nil
}

// remove a directory or file
func (c *PuddleClient) Remove(path string) error {

	// search for path in zookeeper
	exists, _, err := c.zkConn.Exists(c.fsPath + "/" + path)

	// TODO: check if file/dir
	// ed explains how to handle both cases.
	if err != nil {
		return err
	}

	if exists {

		// acquire lock for this path
		// todo: do we have to create lock everytime we want to use it? can we optimize?

		// get the inode from zookeeper
		inode, err := c.getINode(path)

		if err != nil {
			return err
		}

		if inode.IsDir {

			subdirs, _, err := c.zkConn.Children(c.fsPath + "/" + path)

			if err != nil {
				return nil
			}

			// recursively remove subdirectories + files
			c.removeDir(subdirs)

			// once all those are removed, acquire lock and delete this directory

			distLock := CreateDistLock(c.fsPath+"/"+path, c.zkConn)

			distLock.Acquire()

			c.zkConn.Delete(c.fsPath+"/"+path, -1)

			distLock.Release()

		} else {

			// if file, acquire lock and delete

			// TODO: check if file is open?, handle this case here (we already have the lock?, remove from open files?)

			distLock := CreateDistLock(c.fsPath+"/"+path, c.zkConn)

			distLock.Acquire()

			c.zkConn.Delete(c.fsPath+"/"+path, -1)

			distLock.Release()

		}

	} else {
		return err
	}

	return nil
}

// list file & directory names (not full names) under `path`
func (c *PuddleClient) List(path string) ([]string, error) {

	// search for path in zookeeper
	exists, _, err := c.zkConn.Exists(c.fsPath + "/" + path)

	if err != nil {
		return nil, err
	}

	if exists {

		// get the inode from zookeeper
		inode, err := c.getINode(path)

		if err != nil {
			return nil, err
		}

		var output []string

		// if directory, print out subdirectories
		// otherwise, print out file.
		if inode.IsDir {

			// grab children of this directory
			// CONCERN: may return locks on this directory too? do we filter? or should we output?
			output, _, err = c.zkConn.Children(c.fsPath + "/" + path)

			if err != nil {
				return nil, err
			}

		} else {

			// not directory, simply output file name(path/to/file --> file)
			output = append(output, path[strings.LastIndex(path, "/")+1:])
		}

		return output, nil

	} else {
		return nil, err
	}
}

// release zk connection
func (c *PuddleClient) Exit() {
	c.zkConn.Close()
}

// -------------------------- UTILITY/HELPER FUNCTIONS -------------------------- //

// initializes the zookeeper internal file system and locks directory paths
func (c *PuddleClient) initPaths() error {

	// fs path exists
	fsExists, _, err := c.zkConn.Exists(c.fsPath)
	if err != nil {
		return err
	}

	// if fs path does not exist, create it
	if !fsExists {
		_, err = c.zkConn.Create(c.fsPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	// // repeat for locks root
	// locksExists, _, err := c.zkConn.Exists(c.locksPath)
	// if err != nil {
	// 	return err
	// }

	// if !locksExists {
	// 	_, err = c.zkConn.Create(c.locksPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// repeat for tapestry root
	tapestryExists, _, err := c.zkConn.Exists(c.tapestryPath)
	if err != nil {
		return err
	}

	if !tapestryExists {
		_, err = c.zkConn.Create(c.tapestryPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *PuddleClient) findNextFreeFD() int {

	// look through the open files array, find first empty index
	for i, f := range c.openFiles {
		if f == nil {
			return i
		}
	}

	// if no empty index return -1 (MAX OPEN FILES IS 256, SHOULD WE HAVE A LIMIT? TODO: ask about this)
	return -1

}

// helper function to return the tap addr from a client node, eg /tapestry/CLIENTUUID
// func (c *PuddleClient) getFullTapestryAddrPath() string {
// 	return c.tapestryPath + "/" + c.ID
// }

// takes in children of directory we are removing
// if file: acquire lock and remove
// if dir: recur, once done acquire directory lock and remove.
func (c *PuddleClient) removeDir(paths []string) error {

	for _, path := range paths {

		inode, err := c.getINode(path)

		if err != nil {
			return err
		}

		if inode.IsDir {

			// like remove
			// get children if directory, recursively remove all subdirectories
			subdirs, _, err := c.zkConn.Children(c.fsPath + "/" + path)

			if err != nil {
				return nil
			}

			// recursively remove subdirectories + files
			c.removeDir(subdirs)

			// once all those are removed, acquire lock and delete this directory

			distLock := CreateDistLock(c.fsPath+"/"+path, c.zkConn)

			distLock.Acquire()

			c.zkConn.Delete(path, -1)

			distLock.Release()

		} else {

			// if file, acquire lock and delete

			distLock := CreateDistLock(c.fsPath+"/"+path, c.zkConn)

			distLock.Acquire()

			c.zkConn.Delete(c.fsPath+"/"+path, -1)

			distLock.Release()

			return nil
		}
	}
	return nil
}

// helper function given path, decodes byte to return inode.
func (c *PuddleClient) getINode(path string) (*inode, error) {

	// get the inode from zookeeper
	data, _, err := c.zkConn.Get(c.fsPath + "/" + path)

	if err != nil {
		return nil, err
	}

	// unmarshal the inode
	newFileinode, err := decodeInode(data)
	if err != nil {
		return nil, err
	}

	return newFileinode, nil
}

func (c *PuddleClient) getRandomTapestryNode() (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // seed with current time

	// get children of tapestry/node- to get tap nodes
	nodes, _, err := c.zkConn.Children(c.tapestryPath + "/node-")

	if err != nil {
		return "", err
	}

	// select random node to connect to
	selectedNode := nodes[r.Intn(len(nodes))]

	return selectedNode, nil
}
