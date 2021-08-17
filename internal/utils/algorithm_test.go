package utils

import (
	"fmt"
	"testing"
)

func algorithmSlicesEqual(expected [][]Algorithm, received [][]Algorithm) bool {
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

func createTestAlgorithm(init int) Algorithm {
	return Algorithm{
		UserID:      uint64(init),
		Subject:     fmt.Sprintf("Subject%v", init),
		Description: fmt.Sprintf("Description%v", init),
	}
}

func createTestAlgorithmSliceRange(start int, end int) []Algorithm {
	if end < start {
		panic("end is less that start")
	}
	slice := make([]Algorithm, 0, end-start)
	for i := start; i < end; i++ {
		slice = append(slice, createTestAlgorithm(i))
	}
	return slice
}

func createTestAlgorithmSlice(value int) []Algorithm {
	return createTestAlgorithmSliceRange(value, value+1)
}

func TestAlgorithmSplitNil(t *testing.T) {
	slices, err := SplitAlgorithmsToBulks(nil, 1)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}

	var expectedSlices [][]Algorithm = nil
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplitEmpty(t *testing.T) {
	emptySlice := createTestAlgorithmSlice(0)
	slices, err := SplitAlgorithmsToBulks(emptySlice, 1)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}

	expectedSlices := [][]Algorithm{createTestAlgorithmSlice(0)}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit1In1(t *testing.T) {
	oneSlice := createTestAlgorithmSlice(1)
	slices, err := SplitAlgorithmsToBulks(oneSlice, 1)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}
	expectedSlices := [][]Algorithm{createTestAlgorithmSlice(1)}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit1In0(t *testing.T) {
	oneSlice := createTestAlgorithmSlice(1)
	slices, err := SplitAlgorithmsToBulks(oneSlice, 0)
	if err == nil {
		t.Fatalf("error %v retuned does not match expected %v ", err, ZeroOrNegativeChunksSize)
	}
	if slices != nil {
		t.Fatalf("slices is not nil ")
	}
}

func TestAlgorithmSplit2In2(t *testing.T) {
	slice := createTestAlgorithmSliceRange(1, 3)
	slices, err := SplitAlgorithmsToBulks(slice, 2)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]Algorithm{createTestAlgorithmSliceRange(1, 3)}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit3In2(t *testing.T) {
	slice := createTestAlgorithmSliceRange(1, 4)
	slices, err := SplitAlgorithmsToBulks(slice, 2)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRange(1, 3),
		createTestAlgorithmSlice(3)}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit5In2(t *testing.T) {
	slice := createTestAlgorithmSliceRange(1, 6)
	slices, err := SplitAlgorithmsToBulks(slice, 2)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRange(1, 3),
		createTestAlgorithmSliceRange(3, 5),
		createTestAlgorithmSlice(5),
	}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit5In3(t *testing.T) {
	slice := createTestAlgorithmSliceRange(1, 6)
	slices, err := SplitAlgorithmsToBulks(slice, 3)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
		return
	}
	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRange(1, 4),
		createTestAlgorithmSliceRange(4, 6),
	}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit5InMaxInt(t *testing.T) {
	const MaxUint = ^uint(0)
	const MaxInt = int(MaxUint >> 1)

	oneSlice := createTestAlgorithmSliceRange(1, 6)
	slices, err := SplitAlgorithmsToBulks(oneSlice, MaxInt)
	if err != nil {
		t.Fatalf("error %v retuned from SplitToChunksInt ", err)
	}
	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRange(1, 6),
	}
	if !algorithmSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func compareAlgorithmMaps(t *testing.T, expected map[uint64]Algorithm, received map[uint64]Algorithm) {
	if len(received) != len(expected) {
		t.Fatalf("received map size %v doesn't match expeced %v", len(received), len(expected))
		return
	}

	for expectedKey, expectedValue := range expected {
		value, found := received[expectedKey]
		if !found {
			t.Fatalf("expected key[%v] is missing in received map", expectedKey)
			return
		}
		if value != expectedValue {
			t.Fatalf("value %v key[%v] does not match expected %v value", value, expectedKey, expectedValue)
			return
		}
	}
}

func TestAlgorithmNilSlice(t *testing.T) {
	var input []Algorithm
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	if len(reversedMap) != 0 {
		t.Fatalf("size %v is not zero", len(reversedMap))
	}
}

func TestAlgorithmEmptySlice(t *testing.T) {
	input := make([]Algorithm, 0)
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	if len(reversedMap) != 0 {
		t.Fatalf("size %v is not 0", len(reversedMap))
	}
}

func TestAlgorithm1(t *testing.T) {
	input := createTestAlgorithmSlice(1)
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	expectedMap := map[uint64]Algorithm{
		1: createTestAlgorithm(1),
	}
	compareAlgorithmMaps(t, expectedMap, reversedMap)
}

func TestAlgorithm1_2(t *testing.T) {
	input := createTestAlgorithmSliceRange(1, 3)
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	expectedMap := map[uint64]Algorithm{
		1: createTestAlgorithm(1),
		2: createTestAlgorithm(2),
	}
	compareAlgorithmMaps(t, expectedMap, reversedMap)
}

func TestAlgorithm1_1(t *testing.T) {
	input := []Algorithm{createTestAlgorithm(1), createTestAlgorithm(1)}
	reversedMap, err := AlgorithmSliceToMap(input)
	if err == nil {
		t.Fatalf("expected error")
	}
	if reversedMap != nil {
		t.Fatalf("expected nil map")
	}
}
