package main

import "math"

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func clamp(value, inf, sup int) int {
	return min(max(value, inf), sup)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getNearest(value int, slice []int) int {
	i := 0
	for i < len(slice) && value > slice[i] {
		i++
	}
	if i == len(slice) || i > 0 && abs(value-slice[i-1]) < abs(value-slice[i]) {
		i--
	}
	return slice[i]
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func c(n, k int) int {
	if k > n {
		return 0
	}
	return factorial(n) / (factorial(k) * factorial(n-k))
}
func binomialPMF(n, k int, p float64) float64 {
	return float64(c(n, k)) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

func rightTailedBinomialTest(n, k int, p float64) float64 {
	pValue := 0.0
	for i := k; i <= n; i++ {
		pValue += binomialPMF(n, i, p)
	}
	return pValue
}
