package algorithm

import (
	"fmt"
	"github.com/ozonva/ova-algorithm-api/internal/numerics"
	"testing"
)

func algorithmSlicesEqual(expected []Algorithm, got []Algorithm) bool {
	if len(expected) != len(got) {
		return false
	}

	for j := 0; j < len(expected); j++ {
		if expected[j] != got[j] {
			return false
		}
	}

	return true
}

func algorithmSlicesOfSlicesEqual(expected [][]Algorithm, received [][]Algorithm) bool {
	if len(expected) != len(received) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !algorithmSlicesEqual(expected[i], received[i]) {
			return false
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

func createTestAlgorithmSliceRangeInclusive(start, end int) []Algorithm {
	if end < start {
		panic("end is less that start")
	}
	size := end - start + 1
	slice := make([]Algorithm, 0, size)
	for i := start; i <= end; i++ {
		slice = append(slice, createTestAlgorithm(i))
	}
	return slice
}

func createTestAlgorithmSlice(value int) []Algorithm {
	return createTestAlgorithmSliceRangeInclusive(value, value)
}

func TestAlgorithmSplitNil(t *testing.T) {
	slices := SplitAlgorithmsToBulks(nil, 1)

	var expectedSlices [][]Algorithm = nil
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplitEmpty(t *testing.T) {
	emptySlice := createTestAlgorithmSlice(0)
	slices := SplitAlgorithmsToBulks(emptySlice, 1)

	expectedSlices := [][]Algorithm{createTestAlgorithmSlice(0)}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit1In1(t *testing.T) {
	oneSlice := createTestAlgorithmSlice(1)
	slices := SplitAlgorithmsToBulks(oneSlice, 1)

	expectedSlices := [][]Algorithm{createTestAlgorithmSlice(1)}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit1In0(t *testing.T) {
	oneSlice := createTestAlgorithmSlice(1)
	slices := SplitAlgorithmsToBulks(oneSlice, 0)

	expectedSlices := [][]Algorithm{createTestAlgorithmSlice(1)}

	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit2In2(t *testing.T) {
	slice := createTestAlgorithmSliceRangeInclusive(1, 2)
	slices := SplitAlgorithmsToBulks(slice, 2)

	expectedSlices := [][]Algorithm{createTestAlgorithmSliceRangeInclusive(1, 2)}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit3In2(t *testing.T) {
	slice := createTestAlgorithmSliceRangeInclusive(1, 3)
	slices := SplitAlgorithmsToBulks(slice, 2)

	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRangeInclusive(1, 2),
		createTestAlgorithmSlice(3)}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit5In2(t *testing.T) {
	slice := createTestAlgorithmSliceRangeInclusive(1, 5)
	slices := SplitAlgorithmsToBulks(slice, 2)

	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRangeInclusive(1, 2),
		createTestAlgorithmSliceRangeInclusive(3, 4),
		createTestAlgorithmSlice(5),
	}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit5In3(t *testing.T) {
	slice := createTestAlgorithmSliceRangeInclusive(1, 5)
	slices := SplitAlgorithmsToBulks(slice, 3)

	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRangeInclusive(1, 3),
		createTestAlgorithmSliceRangeInclusive(4, 5),
	}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
		t.Fatalf("slices not equal")
	}
}

func TestAlgorithmSplit5InMaxUint(t *testing.T) {
	oneSlice := createTestAlgorithmSliceRangeInclusive(1, 5)
	slices := SplitAlgorithmsToBulks(oneSlice, numerics.MaxUint)

	expectedSlices := [][]Algorithm{
		createTestAlgorithmSliceRangeInclusive(1, 5),
	}
	if !algorithmSlicesOfSlicesEqual(expectedSlices, slices) {
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
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reversedMap) != 0 {
		t.Fatalf("size %v is not zero", len(reversedMap))
	}
}

func TestAlgorithmEmptySlice(t *testing.T) {
	input := make([]Algorithm, 0)
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reversedMap) != 0 {
		t.Fatalf("size %v is not 0", len(reversedMap))
	}
}

func TestAlgorithm1(t *testing.T) {
	input := createTestAlgorithmSlice(1)
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err)
	}
	expectedMap := map[uint64]Algorithm{
		1: createTestAlgorithm(1),
	}
	compareAlgorithmMaps(t, expectedMap, reversedMap)
}

func TestAlgorithm1_2(t *testing.T) {
	input := createTestAlgorithmSliceRangeInclusive(1, 2)
	reversedMap, err := AlgorithmSliceToMap(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err)
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
	if reversedMap != nil {
		t.Fatalf("expected nil map")
	}
	if err == nil {
		t.Fatalf("expected error")
		return
	}
	if err.Error() != "duplicate UserIDs: 1" {
		t.Fatalf("expected error message: %v", err)
	}
}
