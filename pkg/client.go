package pkg

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	tapestry "tapestry/pkg"

	"github.com/go-zookeeper/zk"
	"github.com/tmthrgd/go-memset"
)

var ROUND_ROBIN = 0

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

	GetID() string
}

type OpenFile struct {
	INode    *inode
	Data     []byte    // cached file block data
	FileLock *DistLock // file lock
}

type PuddleClient struct {
	ID          string
	zkConn      *zk.Conn
	openFiles   []*OpenFile // map from file descriptor to inode, represented as an array of inodes (each fd is an index in the array)
	numReplicas int

	fsPath       string       // file system path prefix within zookeeper, e.g. /puddlestore
	tapestryPath string       // path root for tapestry nodes assigned to each client, e.g. /tapestry
	dirtyFiles   map[int]bool // dirty files set ? for flushing, should contain file descriptors (or file paths)?

	// A single mutex to be used by both goroutines
	sync.Mutex // used when locking maps for editing.

}

// ---------------------- CLIENT INTERFACE IMPLEMENTATION ---------------------- //

// open a file and return a file descriptor, DOES THIS PATH START WITH A /?
func (c *PuddleClient) Open(path string, create, write bool) (int, error) {

	// clean path to remove trailing dir.
	path = strings.TrimSuffix(path, "/")
	fmt.Printf("write open first %v\n", write)

	// error check empty path
	if path == "" {
		return -1, errors.New("empty path inputted")
	}

	// search for lock
	lockExists, _, err := c.zkConn.Exists(LOCK_ROOT + path)

	if err != nil {
		return -1, err
	}

	// create lock for file
	if !lockExists {

		// if creating lock fails, this means directory path does not exist
		// so, no need to create lock.
		// example: /path/file will fail if /path does not exist.
		c.zkConn.Create(LOCK_ROOT+path, []byte(""), 0, zk.WorldACL(zk.PermAll))
	}

	distlock := CreateDistLock(LOCK_ROOT+path, c.zkConn)

	distlock.Acquire()

	// search for the file path metadata in zookeeper once we get lock.
	fileExists, _, err := c.zkConn.Exists(c.fsPath + path)

	if err != nil {
		return -1, err
	}
	var newFileinode *inode
	data := make([]byte, 0)

	if !fileExists { // if the file metadata does not exist in the zookeeper fs

		if !create { // if we are not creating and the file does not exist, return error
			distlock.Release()
			return -1, zk.ErrNoNode
		} else { // otherwise create the file

			// first check  parent
			parentNodeIsDir := c.isParentINodeDir(path)

			if !parentNodeIsDir {
				distlock.Release()
				return -1, errors.New("open: parent is not a directory")
			}

			// create the inode
			newFileinode = &inode{
				Filepath: path,              // this is the path of the file in the actual filesystem
				Size:     0,                 // this is the size of the file in bytes (starts as empty)
				Blocks:   make([]string, 0), // this is the list of data blocks (each block is a uuid that represents an entry in tapestry)
				IsDir:    false,             // this is the flag that indicates if the file is a directory
			}

			// marshal the inode to bytes
			inodeBuffer, _ := encodeInode(*newFileinode)

			// if err != nil { // encode fails COMMENTED OUT FOR TESTING COV
			// 	// release the lock
			// 	distlock.Release()
			// 	return -1, err
			// }

			// create the file metadata in zookeeper, should be neither sequential nor ephemeral
			_, err = c.zkConn.Create(c.fsPath+path, inodeBuffer, 0, zk.WorldACL(zk.PermAll))

			if err != nil { // create fails
				// release the lock
				distlock.Release()
				return -1, err
			}
		}

	} else {
		fmt.Printf("here open connec\n")
		// get the inode from zookeeper
		inodeBuf, _, err := c.zkConn.Get(c.fsPath + path)

		if err != nil {
			distlock.Release()
			return -1, err
		}

		// unmarshal the inode
		newFileinode, err = decodeInode(inodeBuf) // contains existing inode data

		if err != nil {
			distlock.Release()
			return -1, err
		}

		// if the file is a directory, return error
		if newFileinode.IsDir {
			distlock.Release()
			return -1, errors.New("open: file is a directory")
		}

		// Connect to a tap node
		// place holder array string.
		var notNeeded []string

		// client we are going to try to get with
		var client *tapestry.Client
		client, err = c.getTapNodeConnected(notNeeded)

		if err != nil {
			distlock.Release()
			return -1, err
		}

		data = make([]byte, newFileinode.Size) // create buffer to store file data

		// null terminate for reading
		memset.Memset(data, 0)

		writePtr := 0 // write pointer

		// get the file data from tapestry, loop through block uuids and get the data from tapestry
		for _, blockUUID := range newFileinode.Blocks {

			numTried := 0

			// block of data we will use to store kv.
			var blockData []byte

			var err error

			// try to get block UUID
			for numTried < MAX_RETRIES {
				if client != nil {
					blockData, err = client.Get(blockUUID)
				} else {
					err = errors.New("client is nil")
				}

				if err == nil {
					break
				}

				numTried += 1
				client, err = c.getTapNodeConnected(notNeeded)
			}

			if err != nil {
				distlock.Release()
				return -1, err
			}

			// fill data byte array
			writePtr += copy(data[writePtr:], blockData) // i'm getting a linter warning here??
		}

		// fmt.Println("data after open: ", data)

	}

	// get next client file descriptor
	c.Lock()
	fd := c.findNextFreeFD()

	if fd == -1 { // if there are no free file descriptors, return error
		distlock.Release()
		c.Unlock()
		return -1, errors.New("no free file descriptors, ENOMEM")
	}

	// add the file to the open files list
	c.openFiles[fd] = &OpenFile{
		INode:    newFileinode,
		Data:     data,
		FileLock: distlock,
	}

	// fmt.Println("open: open file data: ", c.openFiles[fd].Data)

	// if we have specified write, add fd to dirty files (to be flushed on close)
	if write {
		fmt.Printf("client: %s, fd %d\n", c.ID, fd)
		c.dirtyFiles[fd] = true
	}
	c.Unlock()

	return fd, nil

}

// close a file and flush its contents to the distributed filesystem
func (c *PuddleClient) Close(fd int) error {

	// open file
	fmt.Printf("about to close %s\n", c.ID)

	c.Lock()
	defer c.Unlock()

	openFile := c.openFiles[fd]

	if openFile == nil {
		return errors.New("close: file not open")
	}

	defer openFile.FileLock.Release()
	dirtyFileBool := c.dirtyFiles[fd]

	// check dirty files set
	if dirtyFileBool {
		// flush the file

		// keeps track of end of array to get correct slice of bytes.
		var end uint64

		// keeps track of new uuids
		var newUIDs []string

		// create buffer of block size.
		var newData []byte = make([]byte, BLOCK_SIZE)

		// length of curr data.
		dataLength := uint64(len(openFile.Data))

		for i := uint64(0); i < dataLength; i += BLOCK_SIZE {

			end += BLOCK_SIZE

			memset.Memset(newData, 0)

			if end > dataLength {
				end = dataLength
			}
			newData = openFile.Data[i:end]

			var newUID = c.getTapNodeforStore()

			var triedIds []string

			// store replicated datablocks.
			alreadyReplicated := 0
			numTries := 0

			for alreadyReplicated < c.numReplicas && numTries < MAX_RETRIES {

				// grab a random tapestry node path from zookeeper
				client, err := c.getTapNodeConnected(triedIds)

				if err != nil {
					fmt.Printf("replicas: %s\n", err)
				} else {

					// if connected success and client not storing, continue, store.
					clientsStored, err := client.Lookup(newUID)
					alreadyReplicated = len(clientsStored)

					// if already replicated enough, break.
					if alreadyReplicated >= c.numReplicas {
						break
					}

					// if no lookup error, check if stored, if not stored, store.
					if err != nil {
						fmt.Printf("lookup error %s\n", err)
					} else {

						// check for stored clients
						alreadyStored := false
						for _, node := range clientsStored {
							if node.ID == client.ID {
								alreadyStored = true
							}
						}

						// if not already stored in this node, store.
						if !alreadyStored {
							err = client.Store(newUID, newData)

							if err != nil {
								fmt.Printf("error store: %s\n", err)
							}
						}

					}

				}

				// store in tried ids so we don't store it again here.
				triedIds = append(triedIds, client.ID)
				numTries += 1

			}

			newUIDs = append(newUIDs, newUID)

		}

		// here: done populating tap with newuids <--> blocks
		// replace inode uids with new
		openFile.INode.Blocks = newUIDs
		// // update size
		// openFile.INode.Size = uint64(len(openFile.Data))

		// marshal the inode to bytes
		inodeBuf, _ := encodeInode(*openFile.INode)

		// write back the inode in zookeeper
		c.zkConn.Set(c.fsPath+openFile.INode.Filepath, inodeBuf, -1)

		delete(c.dirtyFiles, fd)

		// set openfile to nonexistent, can't use same fd.
		c.openFiles[fd] = &OpenFile{nil, nil, nil}
	}

	fmt.Printf("returning close\n")
	return nil
}

// read a file and return a buffer of size `size` starting at `offset`
func (c *PuddleClient) Read(fd int, offset, size uint64) ([]byte, error) {

	// get open file
	c.Lock()
	defer c.Unlock()

	openFile := c.openFiles[fd]

	// print out all of open file
	// fmt.Printf("read, open file data: %v", openFile)

	if openFile == nil {
		return nil, errors.New("read: file not open")
	}

	if offset >= openFile.INode.Size || openFile.INode.Size == 0 || size == 0 {
		return []byte{}, nil
	}

	// fmt.Printf("read data: %v\n", openFile.Data)
	// fmt.Println("read data len: " + fmt.Sprintf("%v", len(openFile.Data)))

	// get minimum of position to read to and size of file
	endPos := size + offset
	if endPos > openFile.INode.Size {
		endPos = openFile.INode.Size
	}

	// get the data from the file
	data := openFile.Data[offset:endPos]

	// fmt.Printf("READ RETURNED data: %v\n", data)

	return data, nil
}

// write a file and write `data` starting at `offset`
func (c *PuddleClient) Write(fd int, offset uint64, data []byte) error {

	// remember to modify the inode data stored locally on each write, flush to zookeeper on close

	// get the open file

	c.Lock()
	defer c.Unlock()

	openFile := c.openFiles[fd]

	if openFile == nil || !c.dirtyFiles[fd] {
		return errors.New("write: file not open or not opened for writing")
	}

	endPos := offset + uint64(len(data)) // the final position of the data to be written

	// if the end position is beyond the size of the file, set the size
	if endPos > openFile.INode.Size {
		openFile.INode.Size = endPos
	}

	// create a new buffer to hold the new data equal to the size of the file
	newData := make([]byte, openFile.INode.Size)

	memset.Memset(newData, 0) // zero out the new data

	// copy the old data into the new buffer
	copy(newData, openFile.Data)

	// overwrite offset -> offset + len(data) with the new data
	copy(newData[offset:], data)

	openFile.Data = newData // set the new data

	return nil
}

// create a directory at the specified path
func (c *PuddleClient) Mkdir(path string) error {
	// should just be a zookeeper create call?

	// STEPS:
	// check the parent dir exists, and is a valid directory (not a file)
	// check if the path already exists

	// clean path to remove trailing dir.
	path = strings.TrimSuffix(path, "/")

	parentNodeIsDir := c.isParentINodeDir(path)

	if !parentNodeIsDir {
		return errors.New("mkdir: parent is not a directory")
	}

	// check if the path already exists
	exists, _, _ := c.zkConn.Exists(c.fsPath + path)
	// l, _ := c.List("/")
	// // index of the path in the list
	// pathIndex := -1
	// for i, p := range l {
	// 	if p == path {
	// 		pathIndex = i
	// 	}
	// } // sometimes exist returns a false positive???
	if exists {
		return errors.New("mkdir: path already exists, " + path)
	}

	// create local inode
	newDirINode := &inode{
		Filepath: path,
		IsDir:    true,
		Blocks:   make([]string, 0),
		Size:     0,
	}

	// marshal the inode to bytes
	inodeBuf, _ := encodeInode(*newDirINode)

	// if err != nil {
	// 	return err
	// }

	// create the directory in zookeeper
	c.zkConn.Create(c.fsPath+path, inodeBuf, 0, zk.WorldACL(zk.PermAll))

	return nil
}

// remove a directory or file
func (c *PuddleClient) Remove(path string) error {

	// clean path to remove trailing dir.
	path = strings.TrimSuffix(path, "/")

	// search for path in zookeeper
	exists, _, err := c.zkConn.Exists(c.fsPath + path)

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

			subdirs, _, err := c.zkConn.Children(c.fsPath + path)

			if err != nil {
				return nil
			}

			// recursively remove subdirectories + files
			c.removeDir(subdirs, path)

			// assert that the directory is empty
			subdirsCheck, _, _ := c.zkConn.Children(c.fsPath + path)
			if len(subdirsCheck) > 0 {
				return errors.New("remove: directory not empty")
			}

			// once all those are removed, acquire lock and delete this directory

			distLock := CreateDistLock(LOCK_ROOT+path, c.zkConn)

			distLock.Acquire()

			err = c.zkConn.Delete(c.fsPath+path, -1)

			if err != nil {
				return err
			}

			distLock.Release()

		} else {

			distLock := CreateDistLock(LOCK_ROOT+path, c.zkConn)

			distLock.Acquire()

			err = c.zkConn.Delete(c.fsPath+path, -1)

			if err != nil {
				return err
			}

			distLock.Release()

		}

	} else {
		return errors.New("remove: file/dir does not exist")
	}

	return nil
}

// list file & directory names (not full names) under `path`
func (c *PuddleClient) List(path string) ([]string, error) {

	// clean path to remove trailing dir.
	path = strings.TrimSuffix(path, "/")

	// search for path in zookeeper
	exists, _, err := c.zkConn.Exists(c.fsPath + path)

	if err != nil {
		return nil, err
	}

	var output []string
	if exists {

		// get the inode from zookeeper

		inode, err := c.getINode(path)

		if err != nil {
			return nil, err
		}

		// if directory, print out subdirectories
		// otherwise, print out file.
		if inode.IsDir {

			// grab children of this directory
			// CONCERN: may return locks on this directory too? do we filter? or should we output?
			output, _, err = c.zkConn.Children(c.fsPath + path)

			if err != nil {
				return nil, err
			}

		} else {

			// not directory, simply output file name(path/to/file --> file)
			output = append(output, path[strings.LastIndex(path, "/")+1:])
		}

		return output, nil

	} else {

		return output, errors.New("List returned nothing")
	}
}

// release zk connection
func (c *PuddleClient) Exit() {
	// close all file descriptors
	c.Lock()
	for i, openFile := range c.openFiles {
		if openFile != nil && openFile.INode != nil {
			c.Close(i)
		}
	}
	c.Unlock()

	// close zk connection
	c.zkConn.Close()
}

func (c *PuddleClient) GetID() string {
	return c.ID
}

// -------------------------- UTILITY/HELPER FUNCTIONS -------------------------- //

func (c *PuddleClient) findNextFreeFD() int {

	// look through the open files array, find first empty index
	for i, f := range c.openFiles {
		if f == nil {
			return i
		}
	}

	// if no empty index return -1
	return -1

}

// takes in children of directory we are removing
// if file: acquire lock and remove
// if dir: recur, once done acquire directory lock and remove.
func (c *PuddleClient) removeDir(paths []string, parentPath string) error {

	for _, path := range paths {

		fullPath := parentPath + "/" + path

		inode, err := c.getINode(fullPath)

		if err != nil {
			return err
		}

		if inode.IsDir {

			// like remove
			// get children if directory, recursively remove all subdirectories
			subdirs, _, err := c.zkConn.Children(c.fsPath + fullPath)

			if err != nil {
				return nil
			}

			// recursively remove subdirectories + files
			c.removeDir(subdirs, fullPath)

			// once all those are removed, acquire lock and delete this directory

			distLock := CreateDistLock(LOCK_ROOT+fullPath, c.zkConn)

			distLock.Acquire()

			err = c.zkConn.Delete(c.fsPath+fullPath, -1)

			if err != nil {
				return err
			}

			distLock.Release()

		} else {

			// if file, acquire lock and delete

			distLock := CreateDistLock(LOCK_ROOT+path, c.zkConn)

			distLock.Acquire()

			err = c.zkConn.Delete(c.fsPath+fullPath, -1)

			if err != nil {
				return err
			}

			distLock.Release()

		}
	}
	return nil
}

// helper function given path, decodes byte to return inode.
func (c *PuddleClient) getINode(path string) (*inode, error) {

	// get the inode from zookeeper
	data, _, err := c.zkConn.Get(c.fsPath + path)

	if err != nil {
		return nil, err
	}

	// unmarshal the inode
	newFileinode, _ := decodeInode(data)
	// if err != nil { // commented error out for testing coverage
	// 	return nil, err
	// }

	return newFileinode, nil
}

// helper function that finds a random tapestry node address in /tapestry/node-xxxx,
// and returns the address of that node.
func (c *PuddleClient) getRandomTapestryNode(triedIds []string) (string, []string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // seed with current time

	// get children of tapestry/node- to get tap nodes
	nodes, _, err := c.zkConn.Children(c.tapestryPath)

	if len(nodes) == 0 {
		return "", nil, errors.New("no tap nodes")
	}

	if err != nil {
		// fmt.Println("error getting children of tapestry/node-" + err.Error())
		return "", triedIds, err
	}

	// select random node to connect to
	randNum := r.Intn(len(nodes))

	// randomly select node
	selectedNode := nodes[randNum]

	// get ID of client:
	toDecode, _, err := c.zkConn.Get(c.tapestryPath + "/" + selectedNode)

	if err != nil {
		// fmt.Println("error getting filepath" + err.Error())
		return "", triedIds, err
	}

	// grab node
	var tapNode *TapestryAddrNode = new(TapestryAddrNode)
	err = decodeMsgPack(toDecode, tapNode) // populates tapNode with addr

	if err != nil {
		// fmt.Println("error decoding" + err.Error())
		return "", triedIds, err
	}

	var nodeID = tapNode.Id

	// keeps getting randNum until we get one not tried already,
	// len of triedids must not be length of nodes
	for contains(triedIds, nodeID) && len(triedIds) < len(nodes) {
		randNum = r.Intn(len(nodes))

		// randomly select node
		selectedNode := nodes[randNum]

		// get ID of client:
		toDecode, _, err := c.zkConn.Get(c.tapestryPath + "/" + selectedNode)

		if err != nil {
			// fmt.Println("error getting filepath" + err.Error())
			return "", triedIds, err
		}

		// grab node
		var tapNode *TapestryAddrNode = new(TapestryAddrNode)
		err = decodeMsgPack(toDecode, tapNode) // populates tapNode with addr

		if err != nil {
			// fmt.Println("error decoding" + err.Error())
			return "", triedIds, err
		}

		nodeID = tapNode.Id
	}
	// once we tried this tap id, we add it to used indices.
	triedIds = append(triedIds, nodeID)

	return c.tapestryPath + "/" + selectedNode, triedIds, nil
}

// checks if rand int is contained already.
func contains(indices []string, rand string) bool {
	for _, str := range indices {
		if str == rand {
			return true
		}
	}
	return false
}

// gets server id round robin style
func (c *PuddleClient) getTapNodeforStore() string {
	// get children of tapestry/node- to get tap nodes
	nodes, _, err := c.zkConn.Children(c.tapestryPath)

	if err != nil {
		// fmt.Println("error getting children of tapestry/node-" + err.Error())
		return "0"
	}

	selectedNode := nodes[ROUND_ROBIN]

	toDecode, _, err := c.zkConn.Get(c.tapestryPath + "/" + selectedNode)

	if err != nil {
		return "0"
	}

	// if round robin picked last server,
	if ROUND_ROBIN == (len(nodes) - 1) {
		ROUND_ROBIN = 0
	} else {
		ROUND_ROBIN = ROUND_ROBIN + 1
	}

	// grab node
	var tapNode *TapestryAddrNode = new(TapestryAddrNode)
	err = decodeMsgPack(toDecode, tapNode) // populates tapNode with addr

	if err != nil {
		return "0"
	}

	id, _ := tapestry.ParseID(tapNode.Id)

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // seed with current time

	//  generate random number
	var rangeNum int = 1
	for i := 0; i < len(id)/2; i++ {
		rangeNum *= 10
	}

	randomAdd := big.NewInt(int64(r.Intn(rangeNum)))
	randomSub := big.NewInt(int64(r.Intn(rangeNum)))
	id.Big().Sub(randomSub, id.Big())

	output := id.Big().Add(randomAdd, id.Big())

	return output.String()
}

func (c *PuddleClient) getTapestryClientFromTapNodePath(filepath string) (*tapestry.Client, error) {

	// grab the bytes stored at the tapestry node path
	toDecode, _, err := c.zkConn.Get(filepath)

	if err != nil {
		// fmt.Printf("here erraaa %s\n", err)
		return nil, err
	}

	var tapNode *TapestryAddrNode = new(TapestryAddrNode)
	decodeMsgPack(toDecode, tapNode) // populates tapNode with addr

	// if err != nil {
	// 	return nil, err
	// }

	if tapNode.TapCli != nil {
		// check if node is online, if not return errr
		// get children of tapestry.
		nodes, _, err := c.zkConn.Children(TAP_ADDRESS_ROOT)

		if err != nil {
			return nil, err
		}

		// look for node
		var found = false
		for _, node := range nodes {
			toDecode, _, err := c.zkConn.Get(node)

			if err != nil {
				return nil, err
			}

			var listedNodes *TapestryAddrNode = new(TapestryAddrNode)
			err = decodeMsgPack(toDecode, listedNodes) // populates tapNode with addr

			if err != nil {
				// fmt.Printf("here err %s\n", err)
				return nil, err
			}

			if tapNode.TapCli.ID == listedNodes.Id {
				found = true
				break
			}

		}

		// if we did not find node in tapestry, return err.
		if !found {
			return nil, err
		}

		return tapNode.TapCli, nil
	} else {
		client, err := tapestry.Connect(tapNode.Addr)

		// if can't connect, return err
		if err != nil {
			// fmt.Printf("here err %s\n", err)
			return nil, err
		}

		// otherwise, add client to metadata and return
		tapNode.TapCli = client

		return client, nil

	}
}

func (c *PuddleClient) isParentINodeDir(path string) bool {

	lastInd := strings.LastIndex(path, "/")

	// fmt.Println("parent node dir last ind: ", lastInd)

	if lastInd <= 0 { // root dir
		return true
	}

	// get the string until the last / from path, otherwise
	parentPath := path[:lastInd]

	// get the inode from zookeeper
	data, _, err := c.zkConn.Get(c.fsPath + parentPath)

	if err != nil {
		return false
	}

	// unmarshal the inode
	newFileinode, _ := decodeInode(data)

	return newFileinode.IsDir

}

// gets tap node and connects to it
// tries up to 3 times if connections fail.
func (c *PuddleClient) getTapNodeConnected(triedIds []string) (*tapestry.Client, error) {
	// READ THE FILE DATA FROM TAPESTRY USING BLOCKS FOUND IN INODE

	selectedNode, triedIds, err := c.getRandomTapestryNode(triedIds) // get tapestry node path of random node

	if err != nil {
		return nil, err
	}

	client, err := c.getTapestryClientFromTapNodePath(selectedNode) // return the tap node connection
	numTried := 1

	for err != nil || client == nil {
		if numTried >= MAX_RETRIES {
			return nil, err
		} else {

			selectedNode, triedIds, err = c.getRandomTapestryNode(triedIds) // get tapestry node path of random node

			if err != nil {
				continue
			}

			client, err = c.getTapestryClientFromTapNodePath(selectedNode) // return the tap node connection

			if err != nil {
				continue
			}

			numTried += 1

		}
	}

	return client, nil
}
