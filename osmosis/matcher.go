package osmosis

import (
	"errors"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type contentMatcher func(c Content) bool

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

func classifyAndBuildMatcher(value []byte) contentMatcher {
	matcherType, err := jsonparser.GetString(value, "matcherType")

	if err != nil {
		panic(err)
	}

	if strings.EqualFold(matcherType, "allWordsMatcher") {
		return getAllWordsMatcher(value)
	} else if strings.EqualFold(matcherType, "oneWordMatcher") {
		return getOneWordMatcher(value)
	} else if strings.EqualFold(matcherType, "conditionalMatcher") {
		return getConditionalMatcher(value)
	} else if strings.EqualFold(matcherType, "regexMatcher") {
		return getRegexMatcher(value)
	}

	panic(errors.New("Unknown matcher type"))
}

func getRegexMatcher(value []byte) contentMatcher {
	matcher := regexMatcher{}

	if regexExpression, err := jsonparser.GetString(value, "regexExpression"); err != nil {
		panic(err)
	} else {
		matcher.Regex = regexExpression
	}

	return matcher.asContentMatcher()
}

func getOneWordMatcher(value []byte) contentMatcher {
	matcher := containsAtleastOneWordMatcher{}
	matcher.Words = extractWords(value)
	return matcher.asContentMatcher()
}

func getAllWordsMatcher(value []byte) contentMatcher {
	matcher := containsAllWordsMatcher{}
	matcher.Words = extractWords(value)
	return matcher.asContentMatcher()
}

func getConditionalMatcher(value []byte) contentMatcher {
	expressions := []contentMatcher{}
	conditionType, _ := jsonparser.GetString(value, "condition")
	jsonparser.ArrayEach(value, func(parsedVal []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			panic(err)
		}
		expressions = append(expressions, classifyAndBuildMatcher(parsedVal))
	}, "expressions")

	matcher := conditionalMatcher{}
	matcher.Condition = conditionType
	matcher.Expressions = expressions

	return matcher.asContentMatcher()
}

func (cm *conditionalMatcher) asContentMatcher() contentMatcher {
	return func(c Content) bool {
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
	return func(c Content) bool {
		for _, word := range c.Words {
			for _, wrdToMatch := range caowm.Words {
				if strings.EqualFold(word, wrdToMatch) {
					return true
				}
			}
		}
		return false
	}
}

func (cawm *containsAllWordsMatcher) asContentMatcher() contentMatcher {
	return func(c Content) bool {
		for _, wrdToMatch := range cawm.Words {
			isFound := false
			for _, word := range c.Words {
				if strings.EqualFold(word, wrdToMatch) {
					isFound = true
				}
			}
			if !isFound {
				return false
			}
		}
		return true
	}
}

func (mrm *regexMatcher) asContentMatcher() contentMatcher {
	return func(c Content) bool {
		reg, err := regexp.Compile(mrm.Regex)

		if err != nil {
			panic(err)
		}

		return reg.MatchString(c.SanitizedText)
	}
}

func extractWords(value []byte) []string {
	words := []string{}
	if wordList, err := jsonparser.GetString(value, "words"); err != nil {
		panic(err)
	} else {
		for _, word := range strings.Split(wordList, ",") {
			words = append(words, word)
		}
	}

	return words
}
