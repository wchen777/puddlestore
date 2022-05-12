package pkg

// tapestry "tapestry/pkg"

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
	Blocks   []string
	IsDir    bool // determines if inode is a directory or file.
}

// returns the path of the file
func (i inode) GetFilePath() string {
	return i.Filepath
}

// returns the size of the file (in bytes)
func (i inode) GetFileSize() uint64 {
	return i.Size
}

// returns the blocks of the file
func (i inode) GetBlocks() []string {
	return i.Blocks
}
