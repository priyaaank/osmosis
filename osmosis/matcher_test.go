package osmosis

import (
	"strings"
	"testing"
)

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
			},
			"sections":[
				{
                    "contentSelector": {
                        "selectorType": "textBlockSelector",
                        "fromText" : "Invoice Number",
                        "toText": "Tax Amount",
                    },
                    "contentExtractors": [
                        {
                            "extractorType": "regexExtractor",
                            "regex": "Invoice\s+Number:\s+([a-zA-Z0-9]+-[0-9]+-[0-9]+-[0-9]+)",
                            "attributeName": "invoiceNumber",
                            "defaultValue": "NA",
                            "groupNumber": 1
						}
                    ]
                }
			]
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
			},
			"sections":[
				{
                    "contentSelector": {
                        "selectorType": "textBlockSelector",
                        "fromText" : "Invoice Number",
                        "toText": "Tax Amount",
                    },
                    "contentExtractors": [
                        {
                            "extractorType": "regexExtractor",
                            "regex": "Invoice\s+Number:\s+([a-zA-Z0-9]+-[0-9]+-[0-9]+-[0-9]+)",
                            "attributeName": "invoiceNumber",
                            "defaultValue": "NA",
                            "groupNumber": 1
						}
                    ]
                }
			]
		}
    ]
}
`

func TestContainsAtleastOneWordMatcherReturnsTrueWhenEvenOneWordIsPresentInContent(t *testing.T) {
	wordsExpected := []string{"Ola", "XXX"}
	c := content{OriginalText: contentString}
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
	c := content{OriginalText: contentString}
	c.prepare()

	containsWords := containsAtleastOneWordMatcher{
		Words: wordsExpected,
	}

	if containsWords.asContentMatcher()(c) {
		t.Errorf("Expected %s to not contain any of the words %s", c.SanitizedText, wordsExpected)
	}
}

func TestContainsAllWordMatcherReturnsTrueWhenAllWordsArePresentInContent(t *testing.T) {
	wordsExpected := []string{"Ola", "Convenience", "SGST"}
	c := content{OriginalText: contentString}
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
	c := content{OriginalText: contentString}
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
	c := content{OriginalText: contentString}
	c.prepare()

	simpleRegexMatcher := regexMatcher{Regex: regexEpr}

	matcher, _ := simpleRegexMatcher.asContentMatcher()

	if !matcher(c) {
		t.Errorf("Expected %s to match the regex %s", c.SanitizedText, regexEpr)
	}
}

func TestThatRegexMatcherReturnsFalseWhenThereAreNoMatchesInContent(t *testing.T) {
	regexEpr := "Ola\\sFee"
	c := content{OriginalText: contentString}
	c.prepare()

	simpleRegexMatcher := regexMatcher{Regex: regexEpr}

	matcher, _ := simpleRegexMatcher.asContentMatcher()

	if matcher(c) {
		t.Errorf("Expected %s to not match the regex %s", c.SanitizedText, regexEpr)
	}
}

func TestRegexMatcherShouldThrowErrorWhenRegexIsNotValid(t *testing.T) {
	expectedError := "error parsing regexp: missing closing ): `(abc`"
	regexEpr := "(abc"
	c := content{OriginalText: contentString}
	c.prepare()

	simpleRegexMatcher := regexMatcher{Regex: regexEpr}

	_, err := simpleRegexMatcher.asContentMatcher()

	if err == nil || strings.Compare(err.Error(), expectedError) != 0 {
		t.Errorf("Expected error %s to be raised", expectedError)
	}
}

func TestThatComplexMatcherCombinationIsEvaluatedForPositiveMatch(t *testing.T) {
	c := content{OriginalText: contentString}
	c.prepare()
	templates, err := LoadConfig(strings.NewReader(positiveConfig))

	if err != nil {
		t.Errorf("Did not expect error to be returned. But was %s", err.Error())
	}

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
	c := content{OriginalText: contentString}
	c.prepare()
	templates, err := LoadConfig(strings.NewReader(negativeConfig))

	if err != nil {
		t.Errorf("Did not expect error to be returned. But was %s", err.Error())
	}

	oldTemplate := templates["Ola"]

	if oldTemplate.Name != "Ola" {
		t.Errorf("Expected the template list to contain ola template %s but was %s", oldTemplate.Name, "Ola")
	}

	isMatch := oldTemplate.Matcher(c)

	if isMatch {
		t.Errorf("Expected the ola template to not match the content")
	}
}

func TestThatRegexMatcherWillReturnMatchedSelection(t *testing.T) {
	expectedContent := "9.0%"
	positiveSelector := `{
		"selectorType": "regexSelector",
		"regex" : "(Convenience Fee[\w\s\(\)%]+CGST\s+([\d.%]+))",
		"groupNumber": 2
	}`
	c := content{OriginalText: contentString}
	c.prepare()

	selector, _ := classifyAndBuildSelector([]byte(positiveSelector))

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
	c := content{OriginalText: contentString}
	c.prepare()

	selector, _ := classifyAndBuildSelector([]byte(positiveSelector))

	selectedContent := selector(c)

	if strings.Compare(selectedContent.OriginalText, expectedContent) != 0 {
		t.Errorf("Expected selected text [%s] to match [%s]", selectedContent.SanitizedText, expectedContent)
	}
}
