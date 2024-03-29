package main

import (
	"fmt"
	"strings"

	onelog "github.com/francoispqt/onelog"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"
)

// validate validates the payload and rejects objects with labels that contain
// palindromes.
func validate(payload []byte) ([]byte, error) {
	data := gjson.GetBytes(
		payload,
		"request.object.metadata.labels")

	if !data.Exists() {
		logger.Warn("cannot read object name from metadata: accepting request")
		return kubewarden.AcceptRequest()
	}

	labels := data.Map()

	logger.DebugWithFields("validating object", func(e onelog.Entry) {
		namespace := gjson.GetBytes(payload, "request.object.metadata.namespace").String()
		name := gjson.GetBytes(payload, "request.object.metadata.name").String()
		e.String("name", name)
		e.String("namespace", namespace)
	})

	for label := range labels {
		if IsPalindrome(label) {
			logger.DebugWithFields("rejecting object", func(e onelog.Entry) {
				namespace := gjson.GetBytes(payload, "request.object.metadata.namespace").String()
				name := gjson.GetBytes(payload, "request.object.metadata.name").String()
				e.String("name", name)
				e.String("namespace", namespace)
				e.String("label_name", label)
			})

			return kubewarden.RejectRequest(kubewarden.Message(fmt.Sprintf("rejecting palindrome label '%v'", label)), kubewarden.NoCode)
		}

		logger.DebugWithFields("label OK", func(e onelog.Entry) {
			namespace := gjson.GetBytes(payload, "request.object.metadata.namespace").String()
			name := gjson.GetBytes(payload, "request.object.metadata.name").String()
			e.String("name", name)
			e.String("namespace", namespace)
			e.String("label_name", label)
		})
	}

	return kubewarden.AcceptRequest()
}

// validateSettings validates settings for this policy, currently it accepts
// any settings.
func validateSettings(payload []byte) ([]byte, error) {
	return kubewarden.AcceptSettings()
}

// IsPalindrome will return true if the input string is a palindrome and false
// otherwise.
func IsPalindrome(str string) bool {
	lower := strings.ToLower(str)
	runes := []rune(lower)
	right := len(runes) - 1

	for i, r := range runes {
		if r != runes[right] {
			return false
		}

		if i > len(runes)/2 {
			break
		}

		right--
	}

	return true
}
