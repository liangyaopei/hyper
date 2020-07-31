package hyper

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/spaolacci/murmur3"
)

// extractMSB returns the leading bits
func extractMSB(bits uint64, hi uint8, lo uint8) uint64 {
	m := uint64(((1 << (hi - lo)) - 1) << lo)
	return (bits & m) >> lo
}

// getHashValMur2 returns the murmur3 hash
// value of byte array
func getHashValMur2(data []byte) uint64 {
	hashed := murmur3.New64()
	hashed.Write(data)
	return hashed.Sum64()
}

// getZeroBitTail returns the count of
// consecutive tail bit of n
func getZeroBitTail(n uint64) uint8 {
	var cnt uint8 = 0
	for i := 0; i < bitUsed; i++ {
		if n>>i&1 == 0 {
			cnt++
		} else {
			break
		}
	}
	return cnt
}

func uint16ToBytes(num uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, num)
	return data
}

func uint32ToBytes(num uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, num)
	return data
}

func uint64ToBytes(num uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, num)
	return data
}

func str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
