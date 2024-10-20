package helper

import (
	"math"
)

// Reverse an array of T mutating the array
func Reverse[T any](array []T) {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
}

func Chunk[T any](array []T, chunkSize int) [][]T {
	var chunked [][]T

	for i := 0; i < len(array); i += chunkSize {
		end := int(math.Min(float64(i+chunkSize), float64(len(array))))
		chunked = append(chunked, array[i:end])
	}

	return chunked
}
