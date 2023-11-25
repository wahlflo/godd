package internal

import (
	"io"
	"os"
	"time"
)

type outputFile struct {
	chunkBuffer           chan []byte
	pathToDestinationFile string
	done                  chan struct{}
	readingFinished       chan struct{}
	errors                chan error
	byteCounter           *byteCounter
	chunkSize             int
}

func newOutputFile(configuration configuration, chunkBuffer chan []byte, readingFinished chan struct{}, errors chan error, byteCounter *byteCounter) *outputFile {
	return &outputFile{
		chunkBuffer:           chunkBuffer,
		pathToDestinationFile: configuration.GetDestinationPath(),
		done:                  make(chan struct{}, 1),
		readingFinished:       readingFinished,
		errors:                errors,
		byteCounter:           byteCounter,
		chunkSize:             configuration.GetChunkSizeInBytes(),
	}
}

func (output *outputFile) blockUntilWriteComplete() {
	<-output.done
}

func (output *outputFile) startWriting() {
	go func() {
		var err error
		if output.pathToDestinationFile == "" {
			err = output.writeLoop(os.Stdout)
		} else {
			err = output.writeToFile()
		}
		if err != nil {
			output.errors <- err
		}
		output.done <- struct{}{}
	}()
}

func (output *outputFile) writeToFile() error {
	file, err := os.Create(output.pathToDestinationFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return output.writeLoop(file)
}

func (output *outputFile) writeLoop(writer io.Writer) error {
	for {
		select {
		case newChunk := <-output.chunkBuffer:
			for {
				// write chunk to output
				bytesWritten, err := writer.Write(newChunk)
				if err != nil {
					return err
				}
				output.byteCounter.addTransferredBytes(int64(bytesWritten))
				if bytesWritten == len(newChunk) {
					break
				}
				newChunk = newChunk[bytesWritten:]
			}
		case <-time.After(time.Millisecond * 50):
			select {
			case <-output.readingFinished:
				return nil
			default:
				continue
			}
		}
	}
}
