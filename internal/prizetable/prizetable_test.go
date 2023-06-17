package prizetable

import "testing"

func TestGetPrizeOrdinaryTenDollars(t *testing.T) {

	p := GetPrize("Ordinary", 3, false)

	expectedPrize := "$10"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}

func TestGetPrizeOrdinaryFiftyDollars(t *testing.T) {

	p := GetPrize("Ordinary", 4, false)

	expectedPrize := "$50"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}

func TestGetPrizeOrdinaryGroupThree(t *testing.T) {

	p := GetPrize("Ordinary", 5, false)

	expectedPrize := "Group 3"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}

func TestGetPrizeOrdinaryGroupOne(t *testing.T) {

	p := GetPrize("Ordinary", 6, false)

	expectedPrize := "Group 1"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}

func TestGetPrizeOrdinaryTwentyFiveDollars(t *testing.T) {

	p := GetPrize("Ordinary", 3, true)

	expectedPrize := "$25"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}

func TestGetPrizeOrdinaryGroupFour(t *testing.T) {

	p := GetPrize("Ordinary", 4, true)

	expectedPrize := "Group 4"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}

func TestGetPrizeOrdinaryGroupTwo(t *testing.T) {

	p := GetPrize("Ordinary", 5, true)

	expectedPrize := "Group 2"

	if p != expectedPrize {
		t.Errorf("expecting prize: %s, got %s instead", expectedPrize, p)
	}

}
