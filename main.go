package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TotoDraw struct {
	WinningNumbers   []int
	AdditionalNumber int
}

func NewTotoDraw(numbers string, a string) (TotoDraw, error) {
	sortedNumbers, err := convertStringToUniqueSortedNumbers(numbers)

	errorPreText := "NewTotoDraw error:"

	if err != nil {
		return TotoDraw{}, fmt.Errorf("%s %s", errorPreText, err)
	}

	if len(sortedNumbers) != 6 {
		return TotoDraw{}, fmt.Errorf("%s winning numbers should only have a length of 6", errorPreText)
	}

	addNum, err := convertStringToNumber(a)

	if err != nil {
		return TotoDraw{}, fmt.Errorf("unable to convert additional number: %s", a)
	}

	for _, n := range sortedNumbers {
		if n == addNum {
			return TotoDraw{}, fmt.Errorf("duplicate number found in additional number")
		}
	}

	n := TotoDraw{
		WinningNumbers:   sortedNumbers,
		AdditionalNumber: addNum,
	}

	return n, nil
}

func main() {
	fmt.Println("Hello, world.")
	result, _ := convertStringToUniqueSortedNumbers("01 2 3 ")
	fmt.Printf("result %v\n", result)
}

func convertStringToUniqueSortedNumbers(str string) ([]int, error) {
	errorPreText := "convertStringToUniqueSortedNumbers error:"
	trimmed := strings.Trim(str, " ")
	split := strings.Split(trimmed, " ")
	numberMap := make(map[int]int, len(split))

	var numbers []int

	for _, s := range split {
		num, err := convertStringToNumber(s)
		if err != nil {
			return []int{}, fmt.Errorf("%s fail to convert %s, in string: %s", errorPreText, s, split)
		}

		_, ok := numberMap[num]
		if ok {
			return []int{}, fmt.Errorf("%s should not have duplicate numbers", errorPreText)
		}

		numberMap[num] = 1

		numbers = append(numbers, num)
	}

	sort.Ints(numbers)

	return numbers, nil
}

func convertStringToNumber(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, err
}
