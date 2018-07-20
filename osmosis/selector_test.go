package osmosis

import (
	"strings"
	"testing"
)

func TestThatTextBlockSelectorSelectsItCorrectly(t *testing.T) {
	expectedOutput := "Pvt. Ltd. 5th FloorInfotech Center Domlur"
	positiveSelector := `{
		"selectorType": "textBlockSelector",
		"fromText" : "Pvt.",
		"toText": "Bengaluru"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatTextBlockSelectorSelectsWhenItChoosesFirstAndLastWord(t *testing.T) {
	expectedOutput := "ANI Technologies Pvt. Ltd. 5th FloorInfotech Center Domlur Bengaluru Karnataka 560000 Invoice ID 1IE88NHTQ55547 Customer Name Jacob Description Ola Convenience Fee  1IE88NHTQ55547 Convenience Fee Ride Convenience Fee Play Convenience Fee8 CGST 9.0 SGST 9.0 Total Convenience Fee Fare Authorised"
	positiveSelector := `{
		"selectorType": "textBlockSelector",
		"fromText" : "ANI",
		"toText": "Signatory"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatTextBlockSelectorSelectsTillTheEndWhenTillTextIsNotPresent(t *testing.T) {
	expectedOutput := "Total Convenience Fee Fare Authorised Signatory"
	positiveSelector := `{
		"selectorType": "textBlockSelector",
		"fromText" : "Total",
		"toText": "NonPresent"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatTextBlockSelectorSelectsFromBeginningWhenFromTextIsNotPresent(t *testing.T) {
	expectedOutput := "ANI Technologies Pvt. Ltd. 5th Floor"
	positiveSelector := `{
		"selectorType": "textBlockSelector",
		"fromText" : "NotPresent",
		"toText": "Infotech"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatContentSelectorForTextBlockSelectorSelectsFromBeginningToEndwhenBothFromAndTillTextAreNotPresent(t *testing.T) {
	expectedOutput := "ANI Technologies Pvt. Ltd. 5th FloorInfotech Center Domlur Bengaluru Karnataka 560000 Invoice ID 1IE88NHTQ55547 Customer Name Jacob Description Ola Convenience Fee  1IE88NHTQ55547 Convenience Fee Ride Convenience Fee Play Convenience Fee8 CGST 9.0 SGST 9.0 Total Convenience Fee Fare Authorised Signatory"
	positiveSelector := `{
		"selectorType": "textBlockSelector",
		"fromText" : "NotPresent",
		"toText": "AlsoNotPresent"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatTextBlockSelectorSelectsFromBeginningToEndwhenFromAndTillTextAreNotProvidedInConfig(t *testing.T) {
	expectedOutput := "ANI Technologies Pvt. Ltd. 5th FloorInfotech Center Domlur Bengaluru Karnataka 560000 Invoice ID 1IE88NHTQ55547 Customer Name Jacob Description Ola Convenience Fee  1IE88NHTQ55547 Convenience Fee Ride Convenience Fee Play Convenience Fee8 CGST 9.0 SGST 9.0 Total Convenience Fee Fare Authorised Signatory"
	positiveSelector := `{
		"selectorType": "textBlockSelector"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatLineSelectorWillSelectBetweenSpecifiedLines(t *testing.T) {
	expectedOutput := "Customer Name Jacob Description"
	positiveSelector := `{
		"selectorType": "lineNumberSelector",
		"fromLine" : 2,
		"toLine": 3
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatLineSelectorWillSelectTillLastLineWhenToLineIsOutOfBounds(t *testing.T) {
	expectedOutput := "SGST 9.0 Total Convenience Fee Fare Authorised Signatory"
	positiveSelector := `{
		"selectorType": "lineNumberSelector",
		"fromLine" : 8,
		"toLine": 100
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedOutput) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedOutput)
	}
}

func TestThatLineSelectorWillSelectFromFirstLineWhenFromLineIsMissing(t *testing.T) {
	positiveSelector := `{
		"selectorType": "lineNumberSelector",
		"toLine": 100
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.OriginalText, contentString) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, contentString)
	}
}

func TestThatLineSelectorWillSelectFromStartTillLastLineWhenToLineIsMissing(t *testing.T) {
	positiveSelector := `{
		"selectorType": "lineNumberSelector"
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.OriginalText, contentString) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, contentString)
	}
}

func TestThatNestedSelectorCanBeProvidedAsPartOfConfig(t *testing.T) {
	expectedText := "CGST 9.0 SGST 9.0"
	contentSelector := `{
		"selectorType": "textBlockSelector",
		"fromText" : "Ola Convenience Fee",
		"toText": "Total",
		"contentSelector" : {
			"selectorType": "lineNumberSelector",
			"fromLine": 4,
			"toLine": 5
		}
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(contentSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.SanitizedText, expectedText) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedText)
	}
}
