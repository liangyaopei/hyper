// Package hyper takes formulation used in the repo is from wikipedia
// https://en.wikipedia.org/wiki/HyperLogLog
package hyper

import (
	"fmt"
	"math"
	"sync"
)

const (
	bitUsed = 64
	two32   = 1 << 32
)

type hyperLogLog struct {
	registers []uint8 // registers to store value
	m         uint64  // the size of bucket
	p         uint8   // precision

	lock       *sync.RWMutex
	concurrent bool
}

// New returns a new hyperLogLog
// registers size increases proportionally
// to precision
func New(precision uint32, concurrent bool) *hyperLogLog {
	h := &hyperLogLog{
		registers:  make([]uint8, 1<<precision),
		m:          1 << precision,
		p:          uint8(precision),
		concurrent: concurrent,
	}
	if h.concurrent {
		h.lock = &sync.RWMutex{}
	}
	return h
}

// Add add the byte array
func (h *hyperLogLog) Add(data []byte) *hyperLogLog {
	if h.concurrent {
		h.lock.Lock()
		defer h.lock.Unlock()
	}
	hashVal := getHashValMur2(data)
	idx := extractMSB(hashVal, bitUsed, bitUsed-h.p) // {x63,...,x63-p}
	zeroBits := getZeroBitTail(hashVal) + 1
	if h.registers[idx] < zeroBits {
		h.registers[idx] = zeroBits
	}
	return h
}

// Count returns the cardinality estimate.
func (h *hyperLogLog) Count() uint64 {
	if h.concurrent {
		h.lock.RLock()
		defer h.lock.RUnlock()
	}
	est := calculateEstimate(h.registers)
	// E < 5/2*m
	if est <= float64(h.m)*2.5 {
		v := countZeros(h.registers)
		// If V = 0, use the standard HyperLogLog estimator
		if v == 0 {
			return uint64(est)
		}
		// Otherwise, use Linear Counting: E = m log(m/V)
		return uint64(linearCounting(h.m, v))
	}
	if est < two32/30 {
		return uint64(est)
	}
	// for very large cardinalities approaching the limit of the size of the registers
	return uint64(-two32 * math.Log(1-est/two32))
}

// Merge merges other hyperLogLog
func (h *hyperLogLog) Merge(other *hyperLogLog) error {
	if h.p != other.p {
		return fmt.Errorf("precision is not equal. Current is %d, other is %d", h.p, other.p)
	}
	for idx, weight := range other.registers {
		if weight > h.registers[idx] {
			h.registers[idx] = weight
		}
	}
	return nil
}

// M returns the bucket size,
// i.e., the used memory in byte
func (h *hyperLogLog) BucketSize() uint64 {
	if h.concurrent {
		h.lock.RLock()
		defer h.lock.RUnlock()
	}
	return h.m
}

func (h *hyperLogLog) AddString(s string) *hyperLogLog {
	data := str2Bytes(s)
	return h.Add(data)
}

func (h *hyperLogLog) AddUint16(num uint16) *hyperLogLog {
	data := uint16ToBytes(num)
	return h.Add(data)
}

func (h *hyperLogLog) AddUint32(num uint32) *hyperLogLog {
	data := uint32ToBytes(num)
	return h.Add(data)
}

func (h *hyperLogLog) AddUint64(num uint64) *hyperLogLog {
	data := uint64ToBytes(num)
	return h.Add(data)
}

func (h *hyperLogLog) AddBatch(dataArr [][]byte) *hyperLogLog {
	if h.concurrent {
		h.lock.Lock()
		defer h.lock.Unlock()
	}
	for i := 0; i < len(dataArr); i++ {
		data := dataArr[i]
		hashVal := getHashValMur2(data)
		idx := extractMSB(hashVal, bitUsed, bitUsed-h.p) // {x63,...,x63-p}
		zeroBits := getZeroBitTail(hashVal) + 1
		if h.registers[idx] < zeroBits {
			h.registers[idx] = zeroBits
		}
	}
	return h
}

func (h *hyperLogLog) AddStringBatch(s []string) *hyperLogLog {
	dataArr := make([][]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		dataArr = append(dataArr, str2Bytes(s[i]))
	}
	return h.AddBatch(dataArr)
}

func (h *hyperLogLog) AddUint16Batch(nums []uint16) *hyperLogLog {
	dataArr := make([][]byte, 0, len(nums))
	for i := 0; i < len(nums); i++ {
		dataArr = append(dataArr, uint16ToBytes(nums[i]))
	}
	return h.AddBatch(dataArr)
}

func (h *hyperLogLog) AddUint32Batch(nums []uint32) *hyperLogLog {
	dataArr := make([][]byte, 0, len(nums))
	for i := 0; i < len(nums); i++ {
		dataArr = append(dataArr, uint32ToBytes(nums[i]))
	}
	return h.AddBatch(dataArr)
}

func (h *hyperLogLog) AddUint64Batch(nums []uint64) *hyperLogLog {
	dataArr := make([][]byte, 0, len(nums))
	for i := 0; i < len(nums); i++ {
		dataArr = append(dataArr, uint64ToBytes(nums[i]))
	}
	return h.AddBatch(dataArr)
}
