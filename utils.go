package main

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

func Clamp(value, inf, sup int) int {
	return min(max(value, inf), sup)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GetNearest(value int, slice []int) int {
	i := 0
	for i < len(slice) && value > slice[i] {
		i++
	}
	if i == len(slice) || i > 0 && abs(value-slice[i-1]) < abs(value-slice[i]) {
		i--
	}
	return slice[i]
}
