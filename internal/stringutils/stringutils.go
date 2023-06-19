package stringutils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func ConvertStringToUniqueSortedNumbers(str string) ([]int, error) {
	errorPreText := "convertStringToUniqueSortedNumbers error:"
	trimmed := strings.Trim(str, " ")
	split := strings.Split(trimmed, " ")
	numberMap := make(map[int]int, len(split))

	var numbers []int

	for _, s := range split {
		num, err := ConvertStringToNumber(s)
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

func ConvertStringToNumber(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, err
}
