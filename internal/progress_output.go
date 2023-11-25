package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type progressOutput struct {
	bytesWritten                   *byteCounter
	outputProgress                 bool
	chunkBuffer                    chan []byte
	bottleneckStatisticsUpperLimit int
	bottleneckStatisticsLowerLimit int
	bottleneckStatisticsMaxSize    float32
	sizeOfSourceFile               int64
}

func newProgressOutput(bytesWritten *byteCounter, chunkBuffer chan []byte, config configuration) *progressOutput {
	upperLimit := int(float32(config.GetNumberOfBufferedChunks()) * 0.75)
	lowerLimit := int(float32(config.GetNumberOfBufferedChunks()) * 0.25)

	sizeOfSourceFile, err := getSizeOfSourceFile(config.GetSourcePath())
	if err != nil {
		sizeOfSourceFile = 0
	}

	return &progressOutput{
		bytesWritten:                   bytesWritten,
		outputProgress:                 config.GetDestinationPath() != "",
		chunkBuffer:                    chunkBuffer,
		bottleneckStatisticsUpperLimit: upperLimit,
		bottleneckStatisticsLowerLimit: lowerLimit,
		bottleneckStatisticsMaxSize:    float32(config.GetNumberOfBufferedChunks()),
		sizeOfSourceFile:               sizeOfSourceFile,
	}
}

func getSizeOfSourceFile(pathToSourceFile string) (int64, error) {
	if pathToSourceFile == "" {
		return 0, errors.New("input is the standard input")
	}

	// this method allows to determine the file size also of non-regular files like disks,
	// and it works on Unix and Windows
	file, err := os.Open(pathToSourceFile)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	return file.Seek(0, io.SeekEnd)
}

func (receiver *progressOutput) onEventReadIsComplete() {
	if receiver.outputProgress {
		fmt.Println("[+] reading completed")
	}
}

func (receiver *progressOutput) onEventWriteIsComplete() {
	if receiver.outputProgress {
		bytesAlreadyTransferred := receiver.bytesWritten.getValue()
		transferredData := getTransferredDataAsString(bytesAlreadyTransferred)
		fmt.Println("[+] completed   -   " + transferredData + " copied in total")
	}
}

func (receiver *progressOutput) onEventError(err error) {
	fmt.Println("[!] error: " + err.Error())
}

func (receiver *progressOutput) startPrintingProgress() {
	if receiver.outputProgress {
		go receiver.printingProgressLoop()
	}
}

func (receiver *progressOutput) printingProgressLoop() {
	startTime := time.Now()
	time.Sleep(time.Second * 1)
	for {
		bytesAlreadyTransferred := receiver.bytesWritten.getValue()
		duration := time.Now().Sub(startTime)

		transferredData := getTransferredDataAsString(bytesAlreadyTransferred)
		runningSince := formatDuration(duration)
		speed := calculateSpeed(duration, bytesAlreadyTransferred)
		eta := receiver.calculateETA(duration, bytesAlreadyTransferred)

		bufferStatistic := receiver.getBufferStatistic()

		fmt.Println("[+] running since: " + runningSince + "     transferred data: " + transferredData + "     speed: " + speed + "     " + bufferStatistic + eta)

		time.Sleep(time.Second * 10)
	}
}

func (receiver *progressOutput) getBufferStatistic() string {
	chunksInBuffer := len(receiver.chunkBuffer)

	returnString := fmt.Sprintf("buffer: %d%%", int((float32(chunksInBuffer)/receiver.bottleneckStatisticsMaxSize)*100))

	if chunksInBuffer < receiver.bottleneckStatisticsLowerLimit {
		returnString += " (reading is the bottleneck)"
	} else if chunksInBuffer > receiver.bottleneckStatisticsUpperLimit {
		returnString += " (writing is the bottleneck)"
	}
	return returnString
}

func getTransferredDataAsString(bytesAlreadyTransferred int64) string {
	if bytesAlreadyTransferred < 1024 {
		return fmt.Sprint(bytesAlreadyTransferred) + "B"
	}

	bytesAlreadyTransferred /= 1024 // KB
	if bytesAlreadyTransferred < 1024 {
		return fmt.Sprint(bytesAlreadyTransferred) + "KB"
	}

	bytesAlreadyTransferred /= 1024 // MB
	if bytesAlreadyTransferred < 1024 {
		return fmt.Sprint(bytesAlreadyTransferred) + "MB"
	}

	bytesAlreadyTransferred /= 1024 // GB
	if bytesAlreadyTransferred < 1024 {
		return fmt.Sprint(bytesAlreadyTransferred) + "GB"
	}

	bytesAlreadyTransferred /= 1024 // TB
	return fmt.Sprint(bytesAlreadyTransferred) + "TB"
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%dd %02dh %02dm %02ds", days, hours, minutes, seconds)
}

func calculateSpeed(d time.Duration, bytesAlreadyTransferred int64) string {
	bytesPerSecond := float64(bytesAlreadyTransferred) / d.Seconds()

	bytesPerSecond /= 1024 // KB
	bytesPerSecond /= 1024 // MB

	return fmt.Sprintf("%.2fMB/s", bytesPerSecond)
}

func (receiver *progressOutput) calculateETA(d time.Duration, bytesAlreadyTransferred int64) string {
	if receiver.sizeOfSourceFile == 0 {
		return ""
	}

	bytesPerSecond := float64(bytesAlreadyTransferred) / d.Seconds()

	pendingNumberOfBytes := float64(receiver.sizeOfSourceFile - bytesAlreadyTransferred)
	eta := time.Second * time.Duration(pendingNumberOfBytes/bytesPerSecond)
	return "     ETA: " + formatDuration(eta)
}
