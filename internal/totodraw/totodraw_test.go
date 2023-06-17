package totodraw

import "testing"

func TestWinningNumbersContains(t *testing.T) {
	w := WinningNumbers{1, 2, 3, 4, 5, 6}

	expectedElement := 5
	doesContain := w.Contains(expectedElement)
	if !doesContain {
		t.Errorf("expecting element: %d to be found, got %t instead", expectedElement, doesContain)
	}
}

func TestWinningNumbersNotContains(t *testing.T) {
	w := WinningNumbers{1, 2, 3, 4, 5, 6}

	expectedElement := 7
	doesContain := w.Contains(expectedElement)
	if doesContain {
		t.Errorf("expecting element: %d to be not found, got %t instead", expectedElement, doesContain)
	}
}

func TestGetBetTypeOrdinary(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6}

	expectedType := "Ordinary"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}

func TestGetBetTypeSystemSeven(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6, 7}

	expectedType := "System 7"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}

func TestGetBetTypeSystemEight(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6, 7, 8}

	expectedType := "System 8"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}

func TestGetBetTypeSystemNine(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6, 7, 8, 9}

	expectedType := "System 9"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}

func TestGetBetTypeSystemTen(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	expectedType := "System 10"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}

func TestGetBetTypeSystemEleven(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	expectedType := "System 11"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}

func TestGetBetTypeSystemTwelve(t *testing.T) {
	b := Bet{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	expectedType := "System 12"
	actualType := b.GetBetType()
	if actualType != expectedType {
		t.Errorf("expecting type: %s,  got %s instead", expectedType, actualType)
	}
}
