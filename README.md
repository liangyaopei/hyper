[![Go Report Card](https://goreportcard.com/badge/github.com/liangyaopei/hyper)](https://goreportcard.com/report/github.com/liangyaopei/hyper)
[![GoDoc](https://godoc.org/github.com/liangyaopei/hyper?status.svg)](http://godoc.org/github.com/liangyaopei/hyper)

# HyperLogLog
HyperLogLog implementation in Golang, and it is thread-safe can be used concurrently.
The implementation borrows ideas from [wikipedia](https://en.wikipedia.org/wiki/HyperLogLog)  

## Install
```go
go get -u github.com/liangyaopei/hyper
```
## Example
```go
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
```

## Result
For adding 2,8,...524288 to hyperloglog,the result is as follow



| number         | 2    | 8    | 32   | 128      |   512    | 2048     | 8192     | 32768    | 131072   | 524288   |
| -------------- | ---- | ---- | ---- | -------- | :------: | -------- | -------- | -------- | -------- | -------- |
| dirrerent rate | 0.00 | 0.00 | 0.00 | 0.007812 | 0.001953 | 0.001465 | 0.005981 | 0.002197 | 0.005287 | 0.000223 |

