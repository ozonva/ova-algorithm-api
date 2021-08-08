package utils

import (
	"testing"
)

func compareStringIntMaps(t *testing.T, expected map[string]int, received map[string]int) {
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

func TestNilMap(t *testing.T) {
	var input map[int]string
	reversedMap, err := ReverseMapIntString(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	if len(reversedMap) != 0 {
		t.Fatalf("size %v is not zero", len(reversedMap))
	}
}

func TestEmptyMap(t *testing.T) {
	input := make(map[int]string)
	reversedMap, err := ReverseMapIntString(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	if len(reversedMap) != 0 {
		t.Fatalf("size %v is not 0", len(reversedMap))
	}
}

func TestOneElementMap(t *testing.T) {
	input := map[int]string{1: "One"}
	reversedMap, err := ReverseMapIntString(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	expectedMap := map[string]int{"One": 1}
	compareStringIntMaps(t, expectedMap, reversedMap)
}

func TestTwoElementMap(t *testing.T) {
	input := map[int]string{1: "One", 2: "Two"}
	reversedMap, err := ReverseMapIntString(input)
	if err != nil {
		t.Fatalf("unexpected value %v", err.(ReverseMapStringError).value)
	}
	expectedMap := map[string]int{"One": 1, "Two": 2}
	compareStringIntMaps(t, expectedMap, reversedMap)
}

func TestIrreversibleMap(t *testing.T) {
	input := map[int]string{1: "One", 2: "One"}
	reversedMap, err := ReverseMapIntString(input)
	if err == nil {
		t.Fatalf("expected error")
	}

	mapErr := err.(ReverseMapStringError)

	//map traverse order is unspecified, make key always less than second
	if mapErr.key1 > mapErr.key2 {
		mapErr.key1, mapErr.key2 = mapErr.key2, mapErr.key1
	}
	if mapErr.key1 != 1 {
		t.Fatalf("unexpected value of key1 %v, expected 1", mapErr.key1)
	}
	if mapErr.key2 != 2 {
		t.Fatalf("unexpected value of key2 %v, expected 2", mapErr.key1)
	}
	if mapErr.value != "One" {
		t.Fatalf("unexpected value %v of value, expected One", mapErr.value)
	}
	expectedMap := make(map[string]int)
	compareStringIntMaps(t, expectedMap, reversedMap)
}
