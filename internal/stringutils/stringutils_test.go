package stringutils

import "testing"

func TestConvertStringToSortedNumbersSpaceTrim(t *testing.T) {
	input := " 1 2 3 "
	numbers, err := ConvertStringToUniqueSortedNumbers(input)
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

func TestConvertStringToSortedNumbersLeading0(t *testing.T) {
	input := "01 02 03"
	numbers, err := ConvertStringToUniqueSortedNumbers(input)
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

func TestConvertStringToSortedNumbersSorted(t *testing.T) {
	input := "3 2 1"
	numbers, err := ConvertStringToUniqueSortedNumbers(input)
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

func TestConvertStringToNumber(t *testing.T) {
	input := "01"
	number, err := ConvertStringToNumber(input)
	if err != nil {
		t.Errorf("unexpected error converting %s", input)
	}

	expected := 1

	if number != expected {
		t.Errorf("expecting %d got %d instead", expected, number)
	}
}
