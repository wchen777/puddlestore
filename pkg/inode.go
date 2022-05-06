package pkg

import (
	uuid "github.com/google/uuid"
	// tapestry "tapestry/pkg"
)

type inode struct {
	Filepath      string
	FileSizeBytes uint64
	Blocks        []uuid.UUID
}

// returns the path of the file
func (i inode) GetFilePath() string {
	return i.Filepath
}

// returns the size of the file in bytes
func (i inode) GetFileSize() uint64 {
	return i.FileSizeBytes
}

// returns the blocks of the file
func (i inode) GetBlocks() []uuid.UUID {
	return i.Blocks
}
