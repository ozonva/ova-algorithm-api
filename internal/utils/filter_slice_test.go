package utils

import "testing"

func compareSlices(t *testing.T, expected []int, recieved []int) {
	if len(recieved) != len(expected) {
		t.Fatalf("len recived %v does not match expected %v",
			len(recieved), len(expected))
		return
	}

	for i := 0; i < len(recieved); i++ {
		if recieved[i] != expected[i] {
			t.Fatalf("mismach slices on index %v, expected(%v) recieved(%v)",
				i, expected[i], recieved[i])
			return
		}
	}
}

func TestNil(t *testing.T) {
	filtered := FilterWithFixedSlice(nil)
	compareSlices(t, nil, filtered)
}

func TestEmpty(t *testing.T) {
	emptySlice := make([]int, 0)
	filtered := FilterWithFixedSlice(emptySlice)
	expected := make([]int, 0)
	compareSlices(t, expected, filtered)
}

func Test_1(t *testing.T) {
	emptySlice := []int{1}
	filtered := FilterWithFixedSlice(emptySlice)
	expected := []int{1}
	compareSlices(t, expected, filtered)
}

func Test_1_2(t *testing.T) {
	emptySlice := []int{1, 2}
	filtered := FilterWithFixedSlice(emptySlice)
	expected := []int{1}
	compareSlices(t, expected, filtered)
}

func Test_1_2_3(t *testing.T) {
	emptySlice := []int{1, 2, 3}
	filtered := FilterWithFixedSlice(emptySlice)
	expected := []int{1}
	compareSlices(t, expected, filtered)
}
