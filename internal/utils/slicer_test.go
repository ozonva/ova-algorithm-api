package utils

import (
	"testing"

	"github.com/ozonva/ova-algorithm-api/internal/numerics"
)

func nestedSlicesEqual(expected [][]int, received [][]int) bool {
	if len(expected) != len(received) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		subExpected := expected[i]
		subGot := received[i]

		if len(subExpected) != len(subGot) {
			return false
		}

		for j := 0; j < len(subExpected); j++ {
			if subExpected[j] != subGot[j] {
				return false
			}
		}
	}

	return true
}

func TestSplitNil(t *testing.T) {
	slices, err := SplitToChunksInt(nil, 1)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}

	expectedSlices := make([][]int, 0)
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplitEmpty(t *testing.T) {
	emptySlice := make([]int, 0)
	slices, err := SplitToChunksInt(emptySlice, 1)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}

	expectedSlices := make([][]int, 0)
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplit1In1(t *testing.T) {
	oneSlice := []int{1}
	slices, err := SplitToChunksInt(oneSlice, 1)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}
	expectedSlices := [][]int{{1}}
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplit1In0(t *testing.T) {
	oneSlice := []int{1}
	slices, err := SplitToChunksInt(oneSlice, 0)
	if err != ZeroOrNegativeChunksSize {
		t.Fatalf("error %v retuned does not match expected %v ", err, ZeroOrNegativeChunksSize)
	}
	if slices != nil {
		t.Fatalf("slices is not nil ")
	}
}

func TestSplit2In2(t *testing.T) {
	slice := []int{1, 2}
	slices, err := SplitToChunksInt(slice, 2)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]int{{1, 2}}
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplit3In2(t *testing.T) {
	slice := []int{1, 2, 3}
	slices, err := SplitToChunksInt(slice, 2)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]int{{1, 2}, {3}}
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplit5In2(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	slices, err := SplitToChunksInt(slice, 2)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]int{{1, 2}, {3, 4}, {5}}
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplit5In3(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	slices, err := SplitToChunksInt(slice, 3)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]int{{1, 2, 3}, {4, 5}}
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestSplit5InMaxInt(t *testing.T) {
	oneSlice := []int{1, 2, 3, 4, 5}
	slices, err := SplitToChunksInt(oneSlice, numerics.MaxInt)
	if err != NoError {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}
	expectedSlices := [][]int{{1, 2, 3, 4, 5}}
	if !nestedSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}
