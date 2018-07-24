package osmosis

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type contentMatcher func(c content) bool

type containsAtleastOneWordMatcher struct {
	Words []string
}

type containsAllWordsMatcher struct {
	Words []string
}

type regexMatcher struct {
	Regex string
}

type conditionalMatcher struct {
	Condition   string
	Expressions []contentMatcher
}

func classifyAndBuildMatcher(value []byte) (contentMatcher, error) {
	var matcherType string
	var err error

	if matcherType, err = jsonparser.GetString(value, "matcherType"); err != nil {
		return nil, err
	} else if strings.EqualFold(matcherType, "allWordsMatcher") {
		return getAllWordsMatcher(value)
	} else if strings.EqualFold(matcherType, "oneWordMatcher") {
		return getOneWordMatcher(value)
	} else if strings.EqualFold(matcherType, "conditionalMatcher") {
		return getConditionalMatcher(value)
	} else if strings.EqualFold(matcherType, "regexMatcher") {
		return getRegexMatcher(value)
	}

	return nil, fmt.Errorf("ERROR: Unknown matcher type %s", matcherType)
}

func getRegexMatcher(value []byte) (contentMatcher, error) {
	matcher := regexMatcher{}

	if regexExpression, err := jsonparser.GetString(value, "regexExpression"); err != nil {
		return nil, fmt.Errorf("ERROR: Problem building regex matcher. Error is %s", err.Error())
	} else {
		matcher.Regex = regexExpression
	}

	contentMatcherFunc, err := matcher.asContentMatcher()
	if err != nil {
		return nil, err
	}

	return contentMatcherFunc, nil
}

func getOneWordMatcher(value []byte) (contentMatcher, error) {
	var err error
	matcher := containsAtleastOneWordMatcher{}
	if matcher.Words, err = extractWords(value); err != nil {
		return nil, fmt.Errorf("ERROR: Problem building one word matcher. Error is %s", err.Error())
	}
	return matcher.asContentMatcher(), nil
}

func getAllWordsMatcher(value []byte) (contentMatcher, error) {
	var err error
	matcher := containsAllWordsMatcher{}
	if matcher.Words, err = extractWords(value); err != nil {
		return nil, fmt.Errorf("ERROR: Problem building all words matcher. Error is %s", err.Error())
	}
	return matcher.asContentMatcher(), nil
}

func getConditionalMatcher(value []byte) (contentMatcher, error) {
	matcher := conditionalMatcher{}
	expressions := []contentMatcher{}
	conditionType, _ := jsonparser.GetString(value, "condition")
	var parseError error

	jsonparser.ArrayEach(value, func(parsedVal []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			parseError = err
		}
		matcher, err := classifyAndBuildMatcher(parsedVal)

		if err != nil {
			parseError = err
		}
		expressions = append(expressions, matcher)
	}, "expressions")

	if parseError != nil {
		return nil, fmt.Errorf("Could not parse all expressions successfully for conditional matcher, Error is %s", parseError.Error())
	}

	matcher.Condition = conditionType
	matcher.Expressions = expressions

	return matcher.asContentMatcher(), nil
}

func (cm *conditionalMatcher) asContentMatcher() contentMatcher {
	return func(c content) bool {
		var result bool
		isFirst := true

		for _, matcher := range cm.Expressions {

			if isFirst {
				result = matcher(c)
				isFirst = false
				continue
			}

			if strings.EqualFold("and", cm.Condition) {
				result = result && matcher(c)
			} else {
				result = result || matcher(c)
			}
		}

		return result
	}
}

func (caowm *containsAtleastOneWordMatcher) asContentMatcher() contentMatcher {
	return func(c content) bool {
		for _, wrdToMatch := range caowm.Words {
			if strings.Contains(c.OriginalText, wrdToMatch) {
				return true
			}
		}
		return false
	}
}

func (cawm *containsAllWordsMatcher) asContentMatcher() contentMatcher {
	return func(c content) bool {
		for _, wrdToMatch := range cawm.Words {
			isFound := false
			if strings.Contains(c.OriginalText, wrdToMatch) {
				isFound = true
			}
			if !isFound {
				return false
			}
		}
		return true
	}
}

func (mrm *regexMatcher) asContentMatcher() (contentMatcher, error) {
	compiledRegex, err := regexp.Compile(mrm.Regex)

	if err != nil {
		return nil, err
	}

	return func(c content) bool {
		if err != nil {
			log.Printf("ERROR: Error compiling regex %s. Skipping this match. Returning FALSE", mrm.Regex)
			return false
		}

		return compiledRegex.MatchString(c.SanitizedText)
	}, nil
}

func extractWords(value []byte) ([]string, error) {
	words := []string{}
	if wordList, err := jsonparser.GetString(value, "words"); err != nil {
		return nil, err
	} else {
		for _, word := range strings.Split(wordList, ",") {
			words = append(words, word)
		}
	}

	return words, nil
}
