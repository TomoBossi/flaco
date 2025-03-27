package main

import (
	"fmt"
	"math"
	"slices"
)

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

func mean(slice []float64) (float64, error) {
	n := len(slice)
	if n == 0 {
		return 0.0, fmt.Errorf("Error: Cannot find mean of empty slice.")
	}

	sum := 0.0
	for _, value := range slice {
		sum += value
	}
	return sum / float64(n), nil
}

func median(slice []int) (float64, error) {
	n := len(slice)
	if n == 0 {
		return 0.0, fmt.Errorf("Error: Cannot find median of empty slice.")
	}

	slices.Sort(slice)
	if n%2 == 1 {
		return float64(slice[n/2]), nil
	} else {
		return float64(slice[n/2-1]+slice[n/2]) / 2.0, nil
	}
}

func fmtBool(value bool, fmtTrue, fmtFalse string) string {
	if value {
		return fmtTrue
	}
	return fmtFalse
}
