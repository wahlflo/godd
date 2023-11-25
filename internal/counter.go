package internal

import "sync/atomic"

type byteCounter int64

func (c *byteCounter) addTransferredBytes(delta int64) {
	atomic.AddInt64((*int64)(c), delta)
}

func (c *byteCounter) getValue() int64 {
	return atomic.LoadInt64((*int64)(c))
}
