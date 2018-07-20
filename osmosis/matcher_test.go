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
