package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConvertStringToSortedNumbersSpaceTrim(t *testing.T) {
	input := " 1 2 3 "
	numbers, err := convertStringToUniqueSortedNumbers(input)
	if err != nil {
		t.Errorf("unexpected error converting %s", input)
	}

	expected := []int{1, 2, 3}

	for i, num := range expected {
		if numbers[i] != num {
			t.Errorf("was expecting %v but got %v instead", expected, numbers)
		}
	}
}

func TestConvertStringToSortedNumbersLeading0(t *testing.T) {
	input := "01 02 03"
	numbers, err := convertStringToUniqueSortedNumbers(input)
	if err != nil {
		t.Errorf("unexpected error converting %s", input)
	}

	expected := []int{1, 2, 3}

	for i, num := range expected {
		if numbers[i] != num {
			t.Errorf("was expecting %v but got %v instead", expected, numbers)
		}
	}
}

func TestConvertStringToSortedNumbersSorted(t *testing.T) {
	input := "3 2 1"
	numbers, err := convertStringToUniqueSortedNumbers(input)
	if err != nil {
		t.Errorf("unexpected error converting %s", input)
	}

	expected := []int{1, 2, 3}

	for i, num := range expected {
		if numbers[i] != num {
			t.Errorf("was expecting %v but got %v instead", expected, numbers)
		}
	}
}

func TestConvertStringToNumber(t *testing.T) {
	input := "01"
	number, err := convertStringToNumber(input)
	if err != nil {
		t.Errorf("unexpected error converting %s", input)
	}

	expected := 1

	if number != expected {
		t.Errorf("expecting %d got %d instead", expected, number)
	}
}

func TestNewTotoDraw(t *testing.T) {
	inputWinningNumbers := "7 13 18 19 25 29"
	inputAdditionalNumber := "36"

	totoDraw, err := NewTotoDraw(inputWinningNumbers, inputAdditionalNumber)
	if err != nil {
		t.Errorf("error in NewTotoDraw %v", err)
	}

	expectedWinningNumbers := []int{7, 13, 18, 19, 25, 29}

	for i, n := range expectedWinningNumbers {
		if n != totoDraw.WinningNumbers[i] {
			t.Errorf("was expecting %d but got %d instead", n, totoDraw.WinningNumbers[i])
		}
	}

	expectedAdditionalNumber := 36
	if totoDraw.AdditionalNumber != expectedAdditionalNumber {
		t.Errorf("was expecting %d but got %d instead", expectedAdditionalNumber, totoDraw.AdditionalNumber)
	}
}

func TestNewTotoDrawWrongWinningNumberLength(t *testing.T) {
	inputWinningNumbers := "7 13 18 19 25 29 30"
	inputAdditionalNumber := "36"

	_, err := NewTotoDraw(inputWinningNumbers, inputAdditionalNumber)

	expectedErrorString := "NewTotoDraw error: winning numbers should only have a length of 6"
	actualErrorString := fmt.Sprint(err)
	if actualErrorString != expectedErrorString {
		t.Errorf("expected '%s' but got: '%s' instead", expectedErrorString, actualErrorString)
	}
}

func TestNewTotoDrawDuplicateWinningNumber(t *testing.T) {
	inputWinningNumbers := "7 13 18 19 29 29"
	inputAdditionalNumber := "36"

	_, err := NewTotoDraw(inputWinningNumbers, inputAdditionalNumber)

	expectedErrorString := "NewTotoDraw error: convertStringToUniqueSortedNumbers error: should not have duplicate numbers"
	actualErrorString := fmt.Sprint(err)
	if actualErrorString != expectedErrorString {
		t.Errorf("expected '%s' but got: '%s' instead", expectedErrorString, actualErrorString)
	}
}

func TestEndpointNotAllowedMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	expectedStatus := http.StatusMethodNotAllowed
	actualStatus := res.StatusCode
	if actualStatus != expectedStatus {
		t.Errorf("expected %d to be nil got %d", expectedStatus, actualStatus)
	}
}

func TestEndpointInvalidWinningNumbers(t *testing.T) {
	serializedPayload := []byte(`{"winningNumbers": "01 02 03 04 05 06 07", "additionalNumber": "08"}`)
	reader := bytes.NewReader(serializedPayload)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	w := httptest.NewRecorder()
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	expectedStatus := http.StatusBadRequest
	actualStatus := res.StatusCode
	if actualStatus != expectedStatus {
		t.Errorf("expected status: %d got %d", expectedStatus, actualStatus)
	}

	var errorResponseBody ErrorResponseBody
	err := json.NewDecoder(res.Body).Decode(&errorResponseBody)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	expectedMessage := "winning numbers should only contain 6 numbers"
	actualMessage := errorResponseBody.Message
	if actualMessage != expectedMessage {
		t.Errorf("expected %s got %s", expectedMessage, actualMessage)
	}
}

func TestEndpointInvalidCharactersInWinningNumber(t *testing.T) {
	serializedPayload := []byte(`{"winningNumbers": "01 02 03 04 05 06 a7", "additionalNumber": "08"}`)
	reader := bytes.NewReader(serializedPayload)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	w := httptest.NewRecorder()
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	expectedStatus := http.StatusBadRequest
	actualStatus := res.StatusCode
	if actualStatus != expectedStatus {
		t.Errorf("expected status: %d got %d", expectedStatus, actualStatus)
	}

	var errorResponseBody ErrorResponseBody

	if err := json.NewDecoder(res.Body).Decode(&errorResponseBody); err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	expectedMessage := "unable to parse winning numbers"
	actualMessage := errorResponseBody.Message

	if actualMessage != expectedMessage {
		t.Errorf("expected message: %s got %s", expectedMessage, actualMessage)
	}
}
