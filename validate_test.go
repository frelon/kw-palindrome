package main

import (
	"encoding/json"
	"testing"

	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/stretchr/testify/assert"
)

func TestApproval(t *testing.T) {
	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/approve.json",
		nil)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func TestRejection(t *testing.T) {
	payload, err := kubewarden_testing.BuildValidationRequest(
		"test_data/reject.json",
		nil)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_testing.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != false {
		t.Error("Unexpected approval")
	}

	expected_message := "rejecting palindrome labels"
	if response.Message != expected_message {
		t.Errorf("Got '%s' instead of '%s'", response.Message, expected_message)
	}
}

func TestIsPalindrome_Palindromes(t *testing.T) {
	palindromes := []string{
		"",
		"madam",
		"aba",
		"1221",
		"12321",
		"racecar",
		"RaceCar",
		"tattarrattat",
		"TattarrAttat",
		"Föröf",
		"råfår",
		"£££",
	}

	for _, palindrome := range palindromes {
		assert.True(t, IsPalindrome(palindrome), "Expected %v to be a palindrome", palindrome)
	}
}

func TestIsPalindrome_NotPalindromes(t *testing.T) {
	words := []string{
		"abc",
		"123",
		"hello testing",
		"test_data",
		"word",
	}

	for _, word := range words {
		assert.False(t, IsPalindrome(word), "Expected %v to NOT be a palindrome", word)
	}
}
