package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type TotoDraw struct {
	WinningNumbers   []int `json:"winningNumbers"`
	AdditionalNumber int   `json:"additionaNumber"`
}

type Request struct {
	WinningNumbers   string `json:"winningNumbers"`
	AdditionalNumber string `json:"additionalNumber"`
}

type ErrorResponseBody struct {
	Message string `json:"message"`
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

func handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(""))
		return
	}

	var request Request

	err := parseRequestBody(r, &request)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		errorResponseBody := ErrorResponseBody{
			Message: "error parsing request body",
		}
		json.NewEncoder(w).Encode(errorResponseBody)
		return
	}

	winningNumbers, err := convertStringToUniqueSortedNumbers(request.WinningNumbers)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(""))
		return
	}

	if len(winningNumbers) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		errorResponseBody := ErrorResponseBody{
			Message: "winning numbers should only contain 6 numbers",
		}
		json.NewEncoder(w).Encode(errorResponseBody)

	}

}

func parseRequestBody(r *http.Request, data any) error {
	parsed, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		return err
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		return err
	}
	if err := json.Unmarshal(parsed, data); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func main() {
	p := ":8080"
	http.HandleFunc("/", handler)
	fmt.Printf("starting http server\n")
	fmt.Printf("Listening on: %s", p)

	log.Fatal(http.ListenAndServe(p, nil))
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
