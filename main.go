package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aikchun/totoprizecheck/internal/prizetable"
	"github.com/aikchun/totoprizecheck/internal/stringutils"
	"github.com/aikchun/totoprizecheck/internal/totodraw"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

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
	TotoDraw totodraw.TotoDraw    `json:"totoDraw"`
	Results  []totodraw.BetResult `json:"results"`
}

func newTotoDraw(numbers string, a string) (totodraw.TotoDraw, error) {
	var totoDraw totodraw.TotoDraw
	sortedNumbers, err := stringutils.ConvertStringToUniqueSortedNumbers(numbers)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: err.Error(),
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

	d, err := totodraw.NewTotoDraw(sortedNumbers, addNum)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: err.Error(),
		}

		return totoDraw, writeError(errorResponseBody)
	}

	return d, err
}

func mapBetStringsToBets(betStrings []string) ([]totodraw.Bet, error) {
	bets := make([]totodraw.Bet, len(betStrings))

	for i, b := range betStrings {
		bet, err := stringutils.ConvertStringToUniqueSortedNumbers(b)
		if err != nil {
			return []totodraw.Bet{}, err
		}

		bets[i] = bet

	}
	return bets, nil
}

func matchTotoDrawWithBet(t totodraw.TotoDraw, bet totodraw.Bet) totodraw.BetResult {
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

	betType := bet.GetBetType()

	return totodraw.BetResult{
		Numbers:             bet,
		BetType:             betType,
		NumbersMatched:      count,
		HasAdditionalNumber: matchedAdditionalNumber,
		Prize:               prizetable.GetPrize(betType, count, matchedAdditionalNumber),
	}
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

	draw, err := newTotoDraw(request.WinningNumbers, request.AdditionalNumber)
	if err != nil {
		return response, err
	}

	bets, err := mapBetStringsToBets(request.Bets)
	if err != nil {
		errorResponseBody := ErrorResponseBody{
			Status:  400,
			Message: err.Error(),
		}

		return response, writeError(errorResponseBody)
	}

	results := make([]totodraw.BetResult, len(bets))

	for i, bet := range bets {
		betResult := matchTotoDrawWithBet(draw, bet)
		results[i] = betResult
	}

	response.TotoDraw = draw
	response.Results = results
	return response, nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Couldn't find .env")
	}

	e := os.Getenv("ENVIRONMENT")

	if e != "dev" {
		lambda.Start(lambdaHandler)
	}

	p := ":8080"

	if e == "dev" {
		http.HandleFunc("/", handler)
		fmt.Printf("starting http server\n")
		fmt.Printf("listening: %s\n", p)

		log.Fatal(http.ListenAndServe(p, nil))
	}

}
