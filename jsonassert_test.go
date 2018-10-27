package jsonassert_test

import (
	"fmt"
	"testing"

	"github.com/kinbiko/jsonassert"
)

type fakeT struct {
	receivedMessages []string
}

func (ft *fakeT) Errorf(format string, args ...interface{}) {
	ft.receivedMessages = append(ft.receivedMessages, fmt.Sprintf(format, args...))
}

// Should be able to make assertions against the String representation of a
// JSON payload
func TestAssertString(t *testing.T) {
	tt := []struct {
		payload       string
		assertionJSON string
		args          []interface{}
		expAssertions []string
	}{
		{
			// Simple valid check
			payload:       `{"check": "ok"}`,
			assertionJSON: `{"check": "ok"}`,
			expAssertions: []string{},
		},

		{
			// Unparseable payload
			payload:       `Can't parse this`,
			assertionJSON: `{"check": "ok"}`,
			expAssertions: []string{`The given payload is not JSON: "Can't parse this",
nested error is: invalid character 'C' looking for beginning of value`},
		},

		{
			// Unparseable assertion JSON
			payload:       `{"check": "ok"}`,
			assertionJSON: `Can't parse this`,
			expAssertions: []string{`The expected payload is not JSON: "Can't parse this",
nested error is: invalid character 'C' looking for beginning of value`},
		},

		{
			// Mutiple violations, including string formatting
			payload:       `{"check": "nope", "ok": "nah"}`,
			assertionJSON: `{"check": "%s", "ok": "yup"}`,
			args:          []interface{}{"works"},
			expAssertions: []string{
				`Expected key: "check" to have value "works" but was "nope"`,
				`Expected key: "ok" to have value "yup" but was "nah"`,
			},
		},

		{
			// Payload < Assertion JSON
			payload:       `{"ok": "yup"}`,
			assertionJSON: `{"check": "%s", "ok": "yup"}`,
			args:          []interface{}{"works"},
			expAssertions: []string{
				`Expected key "check" to have value "works" but was not present in the payload`,
			},
		},

		{
			// Payload > Assertion JSON
			payload:       `{"check": "works", "ok": "yup"}`,
			assertionJSON: `{"ok": "yup"}`,
			expAssertions: []string{
				`Unexpected key "check" present in the payload`,
			},
		},

		/*
			/*
				{
					payload:      `{"nested": {"check": "ok"}}`,
					assertionJSON: `{"nested": {"check": "%s"}}`,
					args:         []interface{}{"not ok"},
					expAssertions: []string{
						`Expected key: "nested.check" to have value "ok" but was "not ok"`,
					},
				},
		*/
	}
	for _, tc := range tt {
		ft := new(fakeT)
		ja := jsonassert.New(ft)
		ja.AssertString(tc.payload, tc.assertionJSON, tc.args...)

		msgs := ft.receivedMessages
		if exp, got := len(tc.expAssertions), len(msgs); exp != got {
			t.Errorf("Expected %d error messages to be written, but there were %d", exp, got)
			if len(tc.expAssertions) > 0 {
				t.Errorf("Expected the following messages:")
				for _, msg := range tc.expAssertions {
					t.Errorf(" - %s", msg)
				}
			}

			if len(msgs) > 0 {
				t.Errorf("Got the following messages:")
				for _, msg := range msgs {
					t.Errorf(" - %s", msg)
				}
			}
			return //Don't attempt the following assertions

		}

		// The order of the JSON does not matter, so have to do a double subset check
		// Combines the issues in the end in order to make deciphering the test failure easier to parse
		unexpectedAssertions := ""
		for _, got := range msgs {
			found := false
			for _, exp := range tc.expAssertions {
				if got == exp {
					found = true
				}
			}
			if !found {
				if unexpectedAssertions == "" {
					unexpectedAssertions = "Got unexpected assertion failure:"
				}
				unexpectedAssertions += "\n - " + got
			}
		}

		missingAssertions := ""
		for _, got := range tc.expAssertions {
			found := false
			for _, exp := range msgs {
				if got == exp {
					found = true
				}
			}
			if !found {
				if missingAssertions == "" {
					missingAssertions = "\nExpected assertion failure but was not found:"
				}
				missingAssertions += "\n - " + got
			}
		}

		if totalError := unexpectedAssertions + missingAssertions; totalError != "" {
			t.Errorf("Inconsistent assertions:\n%s", totalError)
		}
	}
}