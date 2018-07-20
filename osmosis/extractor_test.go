package osmosis

import (
	"strings"
	"testing"
)

var contentString = `ANI Technologies Pvt. Ltd. 5th Floor,Infotech Center,, Domlur, Bengaluru, Karnataka 560000 Invoice ID 1IE88NHTQ55547
Customer Name Jacob
Description
Ola Convenience Fee - 1IE88NHTQ55547
Convenience Fee (Ride)
Convenience Fee (Play Convenience Fee(8%))
CGST 9.0%
SGST 9.0%
Total Convenience Fee Fare
Authorised Signatory
`

var testConfig = `{
	"templates": [
		{
			"templateName": "Ola",
			"matchers": {
				"matcherType": "conditionalMatcher",
				"condition": "or",
				"expressions": [
					{
						"matcherType":"conditionalMatcher",
						"condition": "and",
						"expressions": [
							{
								"matcherType": "oneWordMatcher",
								"words": "ANI,Simple,Always"
							},
							{
								"matcherType": "regexMatcher",
								"regexExpression": "ANI\\s+Technologies"
							}
						]
					}
				]
			},
			"sections" : [
                {
                    "contentSelector": {
						"selectorType": "lineNumberSelector",
						"fromLine" : 1,
						"toLine": 2
                    },
                    "contentExtractors": [
                        {
							"extractorType": "regexExtractor",
							"regex": "Invoice ID\s+([A-Z0-9]+)",
							"attributeName": "invoiceNumber",
							"defaultValue":"NA",
							"groupNumber":1
						}
                    ]
                }
            ]
		}
    ]
}
`

var positiveConfig = `{
	"templates": [
		{
			"templateName": "Ola",
			"matchers": {
				"matcherType": "conditionalMatcher",
				"condition": "or",
				"expressions": [
					{
						"matcherType":"conditionalMatcher",
						"condition": "and",
						"expressions": [
							{
								"matcherType": "oneWordMatcher",
								"words": "ANI,Simple,Always"
							},
							{
								"matcherType": "regexMatcher",
								"regexExpression": "ANI\\s+Technologies"
							}
						]
					},
					{
						"matcherType":"conditionalMatcher",
						"condition": "and",
						"expressions": [
							{
								"matcherType": "allWordsMatcher",
								"words": "Uber,Invoice"
							},
							{
								"matcherType": "oneWordMatcher",
								"words": "Flight"
							}
						]
					}
				]
			}
		}
    ]
}
`

var negativeConfig = `{
	"templates": [
		{
			"templateName": "Ola",
			"matchers": {
				"matcherType": "conditionalMatcher",
				"condition": "and",
				"expressions": [
					{
						"matcherType": "oneWordMatcher",
						"words": "ANI,Simple,Always"
					},
					{
						"matcherType":"conditionalMatcher",
						"condition": "or",
						"expressions": [
							{
								"matcherType": "regexMatcher",
								"regexExpression": "ANI\\s+Technologies"
							},
							{
								"matcherType": "allWordsMatcher",
								"words": "Uber,Invoice"
							}
						]
					},
					{
						"matcherType": "oneWordMatcher",
						"words": "Flight"
					}
				]
			}
		}
    ]
}
`

func TestThatContentPreparationPopulatesWords(t *testing.T) {
	c := Content{OriginalText: contentString}
	c.prepare()

	if len(c.Words) != 41 {
		t.Errorf("Word count was incorrect. Expected %d but got %d", 32, len(c.Words))
	}
}

func TestThatStringIsSanitizedBeforeItIsConvertedToWords(t *testing.T) {
	c := Content{OriginalText: contentString}
	c.prepare()

	expected := []string{"ANI", "Technologies", "Pvt.", "Ltd.", "5th", "FloorInfotech", "Center", "Domlur", "Bengaluru", "Karnataka", "560000", "Invoice", "ID", "1IE88NHTQ55547", "Customer", "Name", "Jacob", "Description", "Ola", "Convenience", "Fee", "", "1IE88NHTQ55547", "Convenience", "Fee", "Ride", "Convenience", "Fee", "Play", "Convenience", "Fee8", "CGST", "9.0", "SGST", "9.0", "Total", "Convenience", "Fee", "Fare", "Authorised", "Signatory"}

	if len(expected) != len(c.Words) {
		t.Errorf("Expected %s but got %s words", expected, c.Words)
	}

	for index, element := range c.Words {
		if (strings.Compare(element, expected[index])) != 0 {
			t.Errorf("String was not sanitized in content. Expected %s, but got %s at position %d", expected[index], element, index)
		}
	}
}

func TestThatASanitizedCopyOfContentIsStored(t *testing.T) {
	c := Content{OriginalText: contentString}
	c.prepare()

	expected := "ANI Technologies Pvt. Ltd. 5th FloorInfotech Center Domlur Bengaluru Karnataka 560000 Invoice ID 1IE88NHTQ55547 Customer Name Jacob Description Ola Convenience Fee  1IE88NHTQ55547 Convenience Fee Ride Convenience Fee Play Convenience Fee8 CGST 9.0 SGST 9.0 Total Convenience Fee Fare Authorised Signatory"

	if strings.Compare(expected, c.SanitizedText) != 0 {
		t.Errorf("Content not sanitized. Expected %s but got %s", expected, c.SanitizedText)
	}
}

func TestContainsAtleastOneWordMatcherReturnsTrueWhenEvenOneWordIsPresentInContent(t *testing.T) {
	wordsExpected := []string{"ola", "XXX"}
	c := Content{OriginalText: contentString}
	c.prepare()

	containsWords := containsAtleastOneWordMatcher{
		Words: wordsExpected,
	}

	if !containsWords.asContentMatcher()(c) {
		t.Errorf("Expected %s to contain at least one of the words %s", c.SanitizedText, wordsExpected)
	}
}

func TestContainsAtleastOneWordMatcherReturnsFalseWhenNoneOfTheWordsArePresentInContent(t *testing.T) {
	wordsExpected := []string{"Magic", "Close"}
	c := Content{OriginalText: contentString}
	c.prepare()

	containsWords := containsAtleastOneWordMatcher{
		Words: wordsExpected,
	}

	if containsWords.asContentMatcher()(c) {
		t.Errorf("Expected %s to not contain any of the words %s", c.SanitizedText, wordsExpected)
	}
}

func TestContainsAllWordMatcherReturnsTrueWhenAllWordsArePresentInContent(t *testing.T) {
	wordsExpected := []string{"ola", "convenience", "sgst"}
	c := Content{OriginalText: contentString}
	c.prepare()

	containsWords := containsAllWordsMatcher{
		Words: wordsExpected,
	}

	if !containsWords.asContentMatcher()(c) {
		t.Errorf("Expected %s to contain all of the words %s", c.SanitizedText, wordsExpected)
	}
}

func TestContainsAllWordMatcherReturnsFalseWhenEvenOneOfTheWordsIsPresentInContent(t *testing.T) {
	wordsExpected := []string{"ola", "convenience", "sgss"}
	c := Content{OriginalText: contentString}
	c.prepare()

	containsWords := containsAllWordsMatcher{
		Words: wordsExpected,
	}

	if containsWords.asContentMatcher()(c) {
		t.Errorf("Expected %s to not contain one of the words %s", c.SanitizedText, wordsExpected)
	}
}

func TestThatRegexMatcherReturnsTrueWhenThereIsAtleastOneMatchInContent(t *testing.T) {
	regexEpr := "Ola\\s+\\w+\\sFee"
	c := Content{OriginalText: contentString}
	c.prepare()

	simpleRegexMatcher := regexMatcher{Regex: regexEpr}

	if !simpleRegexMatcher.asContentMatcher()(c) {
		t.Errorf("Expected %s to match the regex %s", c.SanitizedText, regexEpr)
	}
}

func TestThatRegexMatcherReturnsFalseWhenThereAreNoMatchesInContent(t *testing.T) {
	regexEpr := "Ola\\sFee"
	c := Content{OriginalText: contentString}
	c.prepare()

	simpleRegexMatcher := regexMatcher{Regex: regexEpr}

	if simpleRegexMatcher.asContentMatcher()(c) {
		t.Errorf("Expected %s to not match the regex %s", c.SanitizedText, regexEpr)
	}
}

func TestRegexMatcherShouldThrowErrorWhenRegexIsNotValid(t *testing.T) {
	expectedError := "error parsing regexp: missing closing ): `(abc`"
	regexEpr := "(abc"
	c := Content{OriginalText: contentString}
	c.prepare()

	defer func() {
		r := recover().(error)
		if strings.Compare(r.Error(), expectedError) != 0 {
			t.Errorf("Expected panic to be thrown for an invalid regex expression. Expected error %s was %s", expectedError, r.Error())
		}
	}()

	simpleRegexMatcher := regexMatcher{Regex: regexEpr}

	if simpleRegexMatcher.asContentMatcher()(c) {
		t.Errorf("Expected %s to not match the regex %s", c.SanitizedText, regexEpr)
	}
}

func TestThatComplexMatcherCombinationIsEvaluatedForPositiveMatch(t *testing.T) {
	c := Content{OriginalText: contentString}
	c.prepare()
	templates := LoadConfig([]byte(positiveConfig))

	oldTemplate := templates["Ola"]

	if oldTemplate.Name != "Ola" {
		t.Errorf("Expected the template list to contain ola template %s but was %s", oldTemplate.Name, "Ola")
	}

	isMatch := oldTemplate.Matcher(c)

	if !isMatch {
		t.Errorf("Expected the ola template to match the content")
	}

}

func TestThatComplexMatcherCombinationIsEvaluatedForNegativeMatch(t *testing.T) {
	c := Content{OriginalText: contentString}
	c.prepare()
	templates := LoadConfig([]byte(negativeConfig))

	oldTemplate := templates["Ola"]

	if oldTemplate.Name != "Ola" {
		t.Errorf("Expected the template list to contain ola template %s but was %s", oldTemplate.Name, "Ola")
	}

	isMatch := oldTemplate.Matcher(c)

	if isMatch {
		t.Errorf("Expected the ola template to not match the content")
	}
}

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

func TestThatRegexMatcherWillReturnMatchedSelection(t *testing.T) {
	expectedContent := "9.0%"
	positiveSelector := `{
		"selectorType": "regexSelector",
		"regex" : "(Convenience Fee[\w\s\(\)%]+CGST\s+([\d.%]+))",
		"groupNumber": 2
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.OriginalText, expectedContent) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedContent)
	}
}

func TestThatRegexMatcherWillReturnEmptyContentWhenNothingMatches(t *testing.T) {
	expectedContent := ""
	positiveSelector := `{
		"selectorType": "regexSelector",
		"regex" : "(Convenience Fee[\w\s\(\)%]Boo+CGST\s+([\d.%]+))",
		"groupNumber": 2
	}`
	c := Content{OriginalText: contentString}
	c.prepare()

	selector := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.OriginalText, expectedContent) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedContent)
	}
}

func TestThatRegexExtractorIsUsedToExtractTheValue(t *testing.T) {
	expectedAttributeName := "invoiceNumber"
	expectedAttributeValue := "1IE88NHTQ55547"
	invoiceExtractor := `{
		"extractorType": "regexExtractor",
		"regex": "Invoice ID\s+([A-Z0-9]+)",
		"attributeName": "invoiceNumber",
		"defaultValue":"NA",
		"groupNumber":1
	}`
	c := Content{OriginalText: "Bengaluru, Karnataka 560000 Invoice ID 1IE88NHTQ55547"}
	c.prepare()

	extractor := classifyAndBuildExtractor([]byte(invoiceExtractor))

	extractedContent := extractor(c)

	if strings.Compare(extractedContent.AttributeName, expectedAttributeName) != 0 {
		t.Errorf("Expected attribute name [%s] to match [%s]", extractedContent.AttributeName, expectedAttributeName)
	}

	if strings.Compare(extractedContent.AttributeValue, expectedAttributeValue) != 0 {
		t.Errorf("Expected attribute name [%s] to match [%s]", extractedContent.AttributeValue, expectedAttributeValue)
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

func TestThatTemplateIsBuiltWithMatcherSelectorAndExtractors(t *testing.T) {
	templates := LoadConfig([]byte(testConfig))

	if len(templates) != 1 {
		t.Errorf("Expected atleast one template to be returned")
	}

	template := templates["Ola"]

	if template.Name != "Ola" {
		t.Errorf("Expected Ola template to be present")
	}

	if template.Matcher == nil {
		t.Errorf("Expected atleast one matcher to be set in template")
	}

	if template.Sections[0].Selector == nil {
		t.Errorf("Expected atleast one selector to be set in template")
	}

	if len(template.Sections[0].Extractors) != 1 {
		t.Errorf("Expected atleast one extractor to be set in template")
	}
}

func TestThatKeyValueMapIsReturnedForTextWhenMatchingTemplateIsPresent(t *testing.T) {
	templates := LoadConfig([]byte(testConfig))
	keyValuePairs := templates.ParseText(contentString)

	if len(keyValuePairs) != 1 {
		t.Errorf("Expected at least one key value pair to be returned")
	}

	if keyValuePairs[0].AttributeName != "invoiceNumber" {
		t.Errorf("Expected attribute key to be invoiceNumber")
	}

	if keyValuePairs[0].AttributeValue != "1IE88NHTQ55547" {
		t.Errorf("Expected attribute value to be 1IE88NHTQ55547")
	}

}
