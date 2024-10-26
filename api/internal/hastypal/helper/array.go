package helper

import (
	"math"
)

type ArrayHelper[T any] struct{}

func NewArrayHelper[T any]() *ArrayHelper[T] {
	return &ArrayHelper[T]{}
}

// Reverse an array of T mutating the array
func (*ArrayHelper[T]) Reverse(array []T) {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
}

func (*ArrayHelper[T]) Chunk(array []T, chunkSize int) [][]T {
	var chunked [][]T

	for i := 0; i < len(array); i += chunkSize {
		end := int(math.Min(float64(i+chunkSize), float64(len(array))))
		chunked = append(chunked, array[i:end])
	}

	return chunked
}
