package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Hello, world.")
	result, _ := ConvertStringToSortedNumbers("01 2 3 ")
	fmt.Printf("result %v\n", result)
}

func ConvertStringToSortedNumbers(str string) ([]int, error) {
	trimmed := strings.Trim(str, " ")
	split := strings.Split(trimmed, " ")
	var numbers []int

	for _, s := range split {
		num, err := strconv.Atoi(s)
		if err != nil {
			return []int{}, errors.New(fmt.Sprintf("error converting %s, in string: %s", s, split))
		}
		numbers = append(numbers, num)

	}
	return numbers, nil

}
