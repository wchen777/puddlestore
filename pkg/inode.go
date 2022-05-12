package pkg

import (
	uuid "github.com/google/uuid"
	// tapestry "tapestry/pkg"
)

// INODES:
/*
	- inodes should be created whenever a file is created.
	- inodes should be stored as metadata in zookeeper
	- under /puddlestore in zookeeper, we have znodes that store inode data for each file path
		- the structure of /puddlestore should mimic the file system structure exactly
*/
type inode struct {
	Filepath string // this is the filepath with respect to the actual file system (not the zookeeper file system)
	Size     uint64
	Blocks   []uuid.UUID
	IsDir    bool // determines if inode is a directory or file.
}
