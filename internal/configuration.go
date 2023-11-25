package internal

type configuration interface {
	GetSourcePath() string
	GetDestinationPath() string
	GetChunkSizeInBytes() int
	GetNumberOfBufferedChunks() int
}
