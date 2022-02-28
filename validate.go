package main

import (
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

	for k := range labels {
		if IsPalindrome(k) {
			logger.DebugWithFields("rejecting object", func(e onelog.Entry) {
				namespace := gjson.GetBytes(payload, "request.object.metadata.namespace").String()
				name := gjson.GetBytes(payload, "request.object.metadata.name").String()
				e.String("name", name)
				e.String("namespace", namespace)
				e.String("label_name", k)
			})

			return kubewarden.RejectRequest(kubewarden.Message("rejecting palindrome labels"), kubewarden.NoCode)
		}

		logger.DebugWithFields("label OK", func(e onelog.Entry) {
			namespace := gjson.GetBytes(payload, "request.object.metadata.namespace").String()
			name := gjson.GetBytes(payload, "request.object.metadata.name").String()
			e.String("name", name)
			e.String("namespace", namespace)
			e.String("label_name", k)
		})
	}

	return kubewarden.AcceptRequest()
}

// IsPalindrome will return true if the input string is a palindrome and false
// otherwise.
func IsPalindrome(str string) bool {
	if len(str) == 0 {
		return false
	}

	lower := strings.ToLower(str)
	runes := []rune(lower)
	i := len(runes) - 1

	for _, r := range runes {
		if r != runes[i] {
			return false
		}

		i--
	}

	return true
}
