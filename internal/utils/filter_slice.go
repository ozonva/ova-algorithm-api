package utils

var fixedSlice []int

func init() {
	fixedSlice = []int{2, 3, 5, 7, 11, 13}
}

func foundInFixedSlice(input int) bool {
	for _, value := range fixedSlice {
		if value == input {
			return true
		}
	}
	return false
}

// FilterWithFixedSlice filters provided slice of int and return a copy slice
// without some values
func FilterWithFixedSlice(toFilter []int) []int {
	var result []int
	for _, value := range toFilter {
		if !foundInFixedSlice(value) {
			result = append(result, value)
		}
	}
	return result
}
