package id

import (
	"strconv"
	"sync/atomic"
)

var autoIncrement uint64 = 0

func Seq() uint64 {
	return atomic.SwapUint64(&autoIncrement, autoIncrement+1)
}

func SeqS() string {
	return strconv.FormatUint(Seq(), 10)
}
