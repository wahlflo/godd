package main

import (
	"flag"
	"fmt"
	"github.com/wahlflo/godd/internal"
	"os"
)

type configuration struct {
	sourcePath             string
	destinationPath        string
	chunkSizeInBytes       int
	numberOfBufferedChunks int
}

func (c *configuration) GetSourcePath() string {
	return c.sourcePath
}

func (c *configuration) GetDestinationPath() string {
	return c.destinationPath
}

func (c *configuration) GetChunkSizeInBytes() int {
	return c.chunkSizeInBytes
}

func (c *configuration) GetNumberOfBufferedChunks() int {
	return c.chunkSizeInBytes
}

func main() {
	helpFlag := flag.Bool("h", false, "show help")
	sourcePath := flag.String("source", "", "path to the source file / disk")
	destinationPath := flag.String("destination", "", "path to the destination file / disk")
	chunkSize := flag.Int("chunk-size", 12, "size of chunks the tool reads / writes in KB")
	bufferSize := flag.Int("buffer-size", 500, "number of chunks which are buffered")

	flag.StringVar(sourcePath, "if", *sourcePath, "path to the source file / disk (short form)")
	flag.StringVar(destinationPath, "of", *destinationPath, "path to the destination file / disk (short form)")

	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	application := internal.NewApplication(&configuration{
		sourcePath:             *sourcePath,
		destinationPath:        *destinationPath,
		chunkSizeInBytes:       *chunkSize * 1024,
		numberOfBufferedChunks: *bufferSize,
	})
	application.Copy()
}

func printHelp() {
	fmt.Println("Usage: godd -source SOURCE -destination DESTINATION [OPTIONS]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
