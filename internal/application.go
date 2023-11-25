package internal

import (
	"os"
	"time"
)

type Application struct {
	configuration configuration
	buffer        chan []byte
}

func NewApplication(configuration configuration) *Application {
	buffer := make(chan []byte, configuration.GetNumberOfBufferedChunks())

	return &Application{
		configuration: configuration,
		buffer:        buffer,
	}
}

func (application *Application) Copy() {
	config := application.configuration
	readingFinished := make(chan struct{}, 1)
	errors := make(chan error, 3)
	counter := byteCounter(0)

	input := newInputFile(config, application.buffer, errors)
	output := newOutputFile(config, application.buffer, readingFinished, errors, &counter)

	progressOutputWriter := newProgressOutput(&counter, application.buffer, config)
	progressOutputWriter.startPrintingProgress()

	go application.checkOnErrorsLoop(progressOutputWriter, errors)

	output.startWriting()
	input.startReading()

	input.blockUntilReadComplete()
	progressOutputWriter.onEventReadIsComplete()
	readingFinished <- struct{}{}

	output.blockUntilWriteComplete()
	progressOutputWriter.onEventWriteIsComplete()

	os.Exit(0)
}

func (application *Application) checkOnErrorsLoop(progressOutputWriter *progressOutput, errors chan error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-errors:
			progressOutputWriter.onEventError(err)
			os.Exit(1)
		case <-ticker.C:
			break
		}
	}
}
