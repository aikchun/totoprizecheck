package main

import (
	"encoding/json"
	"fmt"
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
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	TotoDraw TotoDraw `json:"totoDraw"`
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
	defer r.Body.Close()

	m := r.Method

	if m != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(""))
		return
	}

	var request Request

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "error parsing request body",
		}
		json.NewEncoder(w).Encode(errorResponseBody)
		return
	}

	res, err := lambdaHandler(request)
	if err != nil {
		writeErrorHttp(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func writeErrorHttp(w http.ResponseWriter, err error) {
	var errorResponseBody ErrorResponseBody
	err = json.Unmarshal([]byte(fmt.Sprint(err)), &errorResponseBody)
	if errorResponseBody.Status == http.StatusBadRequest {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponseBody)
	}
}

func writeError(e ErrorResponseBody) error {
	eByte, _ := json.Marshal(e)
	return fmt.Errorf(string(eByte))
}

func lambdaHandler(request Request) (Response, error) {
	var response Response

	winningNumbers, err := convertStringToUniqueSortedNumbers(request.WinningNumbers)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "unable to parse winning numbers",
		}

		return response, writeError(errorResponseBody)
	}

	if len(winningNumbers) != 6 {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "winning numbers should only contain 6 numbers",
		}
		return response, writeError(errorResponseBody)
	}

	additionaNumber, err := convertStringToNumber(request.AdditionalNumber)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "unable to parse additional number",
		}

		return response, writeError(errorResponseBody)
	}

	response.TotoDraw = TotoDraw{
		WinningNumbers:   winningNumbers,
		AdditionalNumber: additionaNumber,
	}

	return response, nil
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
