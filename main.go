package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aikchun/totoprizecheck/internal/prizetable"
	"github.com/aikchun/totoprizecheck/internal/stringutils"
)

type TotoDraw struct {
	WinningNumbers   WinningNumbers `json:"winningNumbers"`
	AdditionalNumber int            `json:"additionalNumber"`
}

func (t TotoDraw) Match(bet Bet) BetResult {
	count := 0
	matchedAdditionalNumber := false
	for _, n := range bet {
		if t.WinningNumbers.Contains(n) {
			count += 1
			continue
		}

		if matchedAdditionalNumber {
			continue
		}

		if n == t.AdditionalNumber {
			matchedAdditionalNumber = true
		}
	}

	betType := getBetType(len(bet))

	return BetResult{
		Numbers:             bet,
		BetType:             betType,
		NumbersMatched:      count,
		HasAdditionalNumber: matchedAdditionalNumber,
		Prize:               prizetable.GetPrize(betType, count, matchedAdditionalNumber),
	}
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
	TotoDraw TotoDraw    `json:"totoDraw"`
	Results  []BetResult `json:"results"`
}

type Bet []int
type WinningNumbers []int

func (w WinningNumbers) Contains(i int) bool {
	for _, n := range w {
		if n == i {
			return true
		}
	}
	return false
}

type BetResult struct {
	Numbers             []int  `json:"numbers"`
	BetType             string `json:"betType"`
	NumbersMatched      int    `json:"numbersMatched"`
	HasAdditionalNumber bool   `json:"hasAdditionalNumber"`
	Prize               string `json:"prize"`
}

func newTotoDraw(numbers string, a string) (TotoDraw, error) {
	var totoDraw TotoDraw
	sortedNumbers, err := stringutils.ConvertStringToUniqueSortedNumbers(numbers)

	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "unable to parse winning numbers",
		}

		return totoDraw, writeError(errorResponseBody)
	}

	if len(sortedNumbers) != 6 {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "winning numbers should only contain 6 numbers",
		}
		return totoDraw, writeError(errorResponseBody)
	}

	addNum, err := stringutils.ConvertStringToNumber(a)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: "unable to parse additional number",
		}

		return totoDraw, writeError(errorResponseBody)
	}

	for _, n := range sortedNumbers {
		if n == addNum {
			errorResponseBody := ErrorResponseBody{
				Status:  400,
				Message: "duplicate number found in additional number",
			}

			return totoDraw, writeError(errorResponseBody)
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

	totoDraw, err := newTotoDraw(request.WinningNumbers, request.AdditionalNumber)
	if err != nil {
		return response, err
	}

	bets, err := convertBetStringsToBets(request.Bets)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: err.Error(),
		}

		return response, writeError(errorResponseBody)
	}

	results := make([]BetResult, len(bets))

	for i, bet := range bets {
		betResult := totoDraw.Match(bet)
		results[i] = betResult
	}

	response.TotoDraw = totoDraw
	response.Results = results
	return response, nil
}

func convertBetStringsToBets(betStrings []string) ([]Bet, error) {
	bets := make([]Bet, len(betStrings))

	for i, b := range betStrings {
		bet, err := stringutils.ConvertStringToUniqueSortedNumbers(b)
		if err != nil {
			return []Bet{}, err
		}

		bets[i] = bet

	}
	return bets, nil
}

func createBetResult(bet Bet, winningNumbers WinningNumbers, additionalNumber int) BetResult {
	count := 0
	matchedAdditionalNumber := false
	for _, n := range bet {
		if winningNumbers.Contains(n) {
			count += 1
			continue
		}

		if matchedAdditionalNumber {
			continue
		}

		if n == additionalNumber {
			matchedAdditionalNumber = true
		}
	}

	betType := getBetType(len(bet))

	return BetResult{
		Numbers:             bet,
		BetType:             betType,
		NumbersMatched:      count,
		HasAdditionalNumber: matchedAdditionalNumber,
		Prize:               prizetable.GetPrize(betType, count, matchedAdditionalNumber),
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

func main() {
	p := ":8080"
	http.HandleFunc("/", handler)
	fmt.Printf("starting http server\n")
	fmt.Printf("Listening on: %s", p)

	log.Fatal(http.ListenAndServe(p, nil))
}
