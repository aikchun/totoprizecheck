package totodraw

import (
	"fmt"
)

type (
	Bet            []int
	WinningNumbers []int
)

func (b Bet) GetBetType() string {
	length := len(b)
	switch length {
	case 6:
		return "Ordinary"
	case 7, 8, 9, 10, 11, 12:
		return fmt.Sprintf("System %d", length)
	}
	return "unknown"
}

func (w WinningNumbers) Contains(i int) bool {
	lo := 0
	hi := len(w)
	for lo < hi {
		mid := ((hi - lo) / 2) + lo
		if w[mid] == i {
			return true
		}
		if i < w[mid] {
			hi = mid
		}
		if i > w[mid] {
			lo = mid + 1
		}

	}
	return false
}

func (w WinningNumbers) IsValid() bool {
	m := make(map[int]bool)
	for _, n := range w {
		if _, found := m[n]; found {
			return false
		} else {
			m[n] = true
		}
	}
	return true
}

type TotoDraw struct {
	WinningNumbers   WinningNumbers `json:"winningNumbers"`
	AdditionalNumber int            `json:"additionalNumber"`
}

type Request struct {
	WinningNumbers   string   `json:"winningNumbers"`
	AdditionalNumber string   `json:"additionalNumber"`
	Bets             []string `json:"bets"`
}

type BetResult struct {
	Numbers             []int  `json:"numbers"`
	BetType             string `json:"betType"`
	NumbersMatched      int    `json:"numbersMatched"`
	HasAdditionalNumber bool   `json:"hasAdditionalNumber"`
	Prize               string `json:"prize"`
}

func NewTotoDraw(w WinningNumbers, a int) (TotoDraw, error) {
	for _, n := range w {
		if n == a {
			return TotoDraw{}, fmt.Errorf("duplicate number found in additional number")
		}
	}

	return TotoDraw{
		WinningNumbers:   w,
		AdditionalNumber: a,
	}, nil
}
