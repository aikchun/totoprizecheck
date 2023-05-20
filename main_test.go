package main

import (
	"testing"
)

func TestConvertStringToSortedNumbers(t *testing.T) {
	input := "01 2 3 "
	numbers, err := ConvertStringToSortedNumbers(input)

	if err != nil {
		t.Errorf("unexpected error converting %s", input)
	}

	expected := []int{1, 2, 3}

	for i, num := range expected {
		if numbers[i] != num {
			t.Errorf("was expecting %v but got %v instead", expected, numbers)
		}

	}
}
