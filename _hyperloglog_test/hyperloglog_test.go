package _hyperloglog_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/liangyaopei/hyper"
)

func TestHyperLogLog_Add(t *testing.T) {
	precision := uint32(12)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 1; i <= 20; i += 2 {
		n := 1 << i
		res := make([]uint32, 0, n)
		h := hyper.New(precision, false)
		for j := 0; j < n; j++ {
			num := rand.Uint32()
			res = append(res, num)
		}
		h.AddUint32Batch(res)
		cnt := h.Count()
		diff := math.Abs(float64(n)-float64(cnt)) / float64(n)
		t.Logf("exact num:%10d,hyperloglog count:%10d,diff:%10f", n, cnt, diff)
	}
}
