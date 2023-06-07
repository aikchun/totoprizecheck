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
	AdditionalNumber int   `json:"additionalNumber"`
}

type Request struct {
	WinningNumbers   string   `json:"winningNumbers"`
	AdditionalNumber string   `json:"additionalNumber"`
	Bets             []string `json:"bets"`
}

type ErrorResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	TotoDraw TotoDraw   `json:"totoDraw"`
	Matches  []BetMatch `json:"bets"`
}

type BetMatch struct {
	Numbers             []int  `json:"numbers"`
	BetType             string `json:"betType"`
	NumbersMatched      int    `json:"numbersMatched"`
	HasAdditionalNumber bool   `json:"hasAdditionalNumber"`
	Prize               string `json:"prize"`
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

	additionalNumber, err := convertStringToNumber(request.AdditionalNumber)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "unable to parse additional number",
		}

		return response, writeError(errorResponseBody)
	}

	response.TotoDraw = TotoDraw{
		WinningNumbers:   winningNumbers,
		AdditionalNumber: additionalNumber,
	}

	matches := make([]BetMatch, len(request.Bets))

	for i, bet := range request.Bets {
		b, err := convertStringToUniqueSortedNumbers(bet)
		if err != nil {
			errorResponseBody := ErrorResponseBody{
				Status:  400,
				Message: fmt.Sprintf("unable to parse %s into numbers", bet),
			}

			return response, writeError(errorResponseBody)
		}

		match := matchBet(b, winningNumbers, additionalNumber)

		matches[i] = match
	}

	response.Matches = matches
	return response, nil
}

func matchBet(bet []int, winningNumbers []int, additionalNumber int) BetMatch {
	count := 0
	matchedAdditionalNumber := false
	for _, n := range bet {
		for _, m := range winningNumbers {
			if n == m {
				count += 1
			}
		}

		if !matchedAdditionalNumber {
			if n == additionalNumber {
				matchedAdditionalNumber = true
			}
		}
	}

	betType := getBetType(len(bet))

	return BetMatch{
		Numbers:             bet,
		BetType:             betType,
		NumbersMatched:      count,
		HasAdditionalNumber: matchedAdditionalNumber,
		Prize:               calculatePrize(betType, count, matchedAdditionalNumber),
	}
}

func getBetType(length int) string {
	switch length {
	case 6:
		return "Ordinary"
	case 7, 8, 9, 10, 11, 12:
		return fmt.Sprintf("System %d", length)
	}
	return "unknown"
}

func calculatePrize(betType string, numbersMatched int, hasAdditionalNumber bool) string {
	switch betType {
	case "Ordinary":
		return calculateOrdinaryPrize(numbersMatched, hasAdditionalNumber)
	case "System 7":
		return calculateSystemSevenPrize(numbersMatched, hasAdditionalNumber)
	case "System 8":
		return calculateSystemEightPrize(numbersMatched, hasAdditionalNumber)
	case "System 9":
		return calculateSystemNinePrize(numbersMatched, hasAdditionalNumber)
	case "System 10":
		return calculateSystemTenPrize(numbersMatched, hasAdditionalNumber)
	case "System 11":
		return calculateSystemElevenPrize(numbersMatched, hasAdditionalNumber)
	case "System 12":
		return calculateSystemTwelvePrize(numbersMatched, hasAdditionalNumber)
	}
	return ""
}

func calculateOrdinaryPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$10"
		case 4:
			return "$50"
		case 5:
			return "Group 3"
		case 6:
			return "Group 1"

		}
	}

	switch numbersMatched {
	case 3:
		return "$25"
	case 4:
		return "Group 4"
	case 5:
		return "Group 2"
	}

	return "unknown"
}

func calculateSystemSevenPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$40"
		case 4:
			return "$190"
		case 5:
			return "Group 3 + $250"
		case 6:
			return "Group 1 + 3"
		}
	}

	switch numbersMatched {
	case 3:
		return "$85"
	case 4:
		return "Group 4 + $150"
	case 5:
		return "Group 2 + 3 + 4"
	case 6:
		return "Group 1 + 2"
	}

	return "unknown"
}

func calculateSystemEightPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$100"
		case 4:
			return "$460"
		case 5:
			return "Group 3 + $850"
		case 6:
			return "Group 1 + 3 + $750"
		}
	}

	switch numbersMatched {
	case 3:
		return "$190"
	case 4:
		return "Group 4 + $490"
	case 5:
		return "Group 2 + 3 + 4 + $500"
	case 6:
		return "Group 1 + 2 + 3 + 4"
	}

	return "unknown"
}

func calculateSystemNinePrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$200"
		case 4:
			return "$900"
		case 5:
			return "Group 3 + $1,900"
		case 6:
			return "Group 1 + 3 + $2,450"
		}
	}

	switch numbersMatched {
	case 3:
		return "$350"
	case 4:
		return "Group 4 + $1,060"
	case 5:
		return "Group 2 + 3 + 4 + $1,600"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $1,250"
	}

	return "unknown"
}

func calculateSystemTenPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$350"
		case 4:
			return "$1,550"
		case 5:
			return "Group 3 + $3,500"
		case 6:
			return "Group 1 + 3 + $5,300"
		}
	}

	switch numbersMatched {
	case 3:
		return "$575"
	case 4:
		return "Group 4 + $1,900"
	case 5:
		return "Group 2 + 3 + 4 + $3,400"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $3,950"
	}

	return "unknown"
}

func calculateSystemElevenPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$560"
		case 4:
			return "$2,450"
		case 5:
			return "Group 3 + $5,750"
		case 6:
			return "Group 1 + 3 + $9,500"
		}
	}

	switch numbersMatched {
	case 3:
		return "$875"
	case 4:
		return "Group 4 + $3,050"
	case 5:
		return "Group 2 + 3 + 4 + $6,000"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $8,300"
	}

	return "unknown"
}

func calculateSystemTwelvePrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$840"
		case 4:
			return "$3,640"
		case 5:
			return "Group 3 + $8,750"
		case 6:
			return "Group 1 + 3 + $15,250"
		}
	}

	switch numbersMatched {
	case 3:
		return "$1,260"
	case 4:
		return "Group 4 + $4,550"
	case 5:
		return "Group 2 + 3 + 4 + $9,500"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $14,500"
	}

	return "unknown"
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
