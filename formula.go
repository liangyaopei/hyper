// Package hyper takes formulation used in the repo is from wikipedia
// https://en.wikipedia.org/wiki/HyperLogLog
package hyper

import "math"

// calculateEstimate returns
// E = alphaM * m^2 * Z
// Z = sum(2^(-M[j]))
func calculateEstimate(s []uint8) float64 {
	sum := 0.0
	for _, val := range s {
		sum += 1.0 / float64(uint64(1)<<val)
	}

	m := uint32(len(s))
	fm := float64(m)
	return alphaM(m) * fm * fm / sum
}

// alphaM returns the approximated constant
// for Count() calculation
func alphaM(m uint32) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/float64(m))
	}
}

// linearCounting E = m log(m/V)
// V is the count of registers equal to 0.
func linearCounting(m uint64, v uint64) float64 {
	fm := float64(m)
	return fm * math.Log(fm/float64(v))
}

// countZeros returns the count of registers equal to 0.
func countZeros(s []uint8) uint64 {
	var c uint64
	for _, v := range s {
		if v == 0 {
			c++
		}
	}
	return c
}
