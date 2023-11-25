package internal

import (
	"io"
	"os"
)

type inputFile struct {
	buffer          chan []byte
	pathToInputFile string
	done            chan struct{}
	chunkSize       int
	errors          chan error
}

func newInputFile(config configuration, buffer chan []byte, errors chan error) *inputFile {
	return &inputFile{
		buffer:          buffer,
		pathToInputFile: config.GetSourcePath(),
		done:            make(chan struct{}, 1),
		chunkSize:       config.GetChunkSizeInBytes(),
		errors:          errors,
	}
}

func (input *inputFile) blockUntilReadComplete() {
	<-input.done
}

func (input *inputFile) startReading() {
	go func() {
		var err error
		if input.pathToInputFile == "" {
			err = input.read(os.Stdin)
		} else {
			err = input.readFromFile()
		}
		if err != nil {
			input.errors <- err
		}
		input.done <- struct{}{}
	}()
}

func (input *inputFile) readFromFile() error {
	file, err := os.Open(input.pathToInputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return input.read(file)
}

func (input *inputFile) read(reader io.Reader) error {
	for {
		chunk := make([]byte, input.chunkSize)
		bytesRead, err := reader.Read(chunk)
		if err != nil {
			if err == io.EOF {
				if bytesRead > 0 {
					input.buffer <- chunk[:bytesRead]
				}
				break
			}
			return err
		}
		input.buffer <- chunk[:bytesRead]
	}
	return nil
}
