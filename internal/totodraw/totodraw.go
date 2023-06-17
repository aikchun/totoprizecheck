package totodraw

import (
	"fmt"

	"github.com/aikchun/totoprizecheck/internal/prizetable"
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

func getBetType(length int) string {
	switch length {
	case 6:
		return "Ordinary"
	case 7, 8, 9, 10, 11, 12:
		return fmt.Sprintf("System %d", length)
	}
	return "unknown"
}
