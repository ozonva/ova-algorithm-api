package utils

import "errors"

// ReverseMapStringError error type for ReverseMapIntString function.
// ok is set if duplicate is found
// key1 the first conflicting key
// key2 the second conflicting key
// value of duplicate sting
type ReverseMapStringError struct {
	error
	key1  int
	key2  int
	value string
}

// ReverseMapIntString converts map[int]string to map[string]int
// so key becomes value and value becomes key. For irreversible map
// where same value exists for two keys error is returned with
// conflict details
func ReverseMapIntString(inputMap map[int]string) (map[string]int, error) {
	reversedMap := make(map[string]int, len(inputMap))
	for inputKey, inputValue := range inputMap {
		reversedValue, found := reversedMap[inputValue]
		if !found {
			reversedMap[inputValue] = inputKey
		} else {
			return nil, ReverseMapStringError{errors.New("duplicate values"), inputKey, reversedValue, inputValue}
		}
	}
	return reversedMap, nil
}
