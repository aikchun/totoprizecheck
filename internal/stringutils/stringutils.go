package stringutils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func ConvertStringToUniqueSortedNumbers(str string) ([]int, error) {
	trimmed := strings.Trim(str, " ")
	split := strings.Split(trimmed, " ")
	numberMap := make(map[int]int, len(split))

	var numbers []int

	for _, s := range split {
		num, err := ConvertStringToNumber(s)
		if err != nil {
			return []int{}, fmt.Errorf("%s: %s", err.Error(), split)
		}

		_, ok := numberMap[num]
		if ok {
			return []int{}, fmt.Errorf("duplicate numbers found: %s", split)
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
		return 0, fmt.Errorf("failed to convert %s", str)
	}

	if num < 1 && num > 49 {
		return 0, fmt.Errorf("number not within range: %s", str)
	}

	return num, err
}
