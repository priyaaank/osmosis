package osmosis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type contentSelector func(c content) content

type textBlockSelector struct {
	FromText string
	ToText   string
}

type lineNumberSelector struct {
	FromLine int64
	ToLine   int64
}

type regexSelector struct {
	RegexPattern string
	GroupNumber  int64
}

func classifyAndBuildSelector(value []byte) (contentSelector, error) {
	var selector contentSelector
	var err error
	var selectorType string

	selectorType, err = jsonparser.GetString(value, "selectorType")

	defer func() (contentSelector, error) {
		if r := recover(); r != nil {
			return nil, fmt.Errorf("ERROR: Error creating a selectorType of %s. Error is %v", selectorType, r)
		}
		return selector, nil
	}()

	if err != nil {
		return nil, fmt.Errorf("ERROR: Could not find tag selectorType in config. Error is " + err.Error())
	}

	if strings.EqualFold(selectorType, "textBlockSelector") {
		selector, err = getTextBlockSelector(value).asContentSelector()
	} else if strings.EqualFold(selectorType, "lineNumberSelector") {
		selector, err = getLineNumberSelector(value).asContentSelector()
	} else if strings.EqualFold(selectorType, "regexSelector") {
		selector, err = getRegexSelector(value).asContentSelector()
	}

	contentSelectorValue, _, _, err := jsonparser.Get(value, "contentSelector")

	if err != nil {
		return selector, nil
	}

	if nestedSelector, err := classifyAndBuildSelector(contentSelectorValue); err != nil {
		return nil, err
	} else if nestedSelector != nil {
		selector = selector.addNestedSelector(nestedSelector)
	}

	return selector, nil
}

func getRegexSelector(value []byte) regexSelector {
	regex, _, _, _ := jsonparser.Get(value, "regex")
	groupNumber, _ := jsonparser.GetInt(value, "groupNumber")

	return regexSelector{
		RegexPattern: string(regex),
		GroupNumber:  groupNumber,
	}
}

func getTextBlockSelector(value []byte) textBlockSelector {
	fromText, _ := jsonparser.GetString(value, "fromText")
	toText, _ := jsonparser.GetString(value, "toText")

	return textBlockSelector{
		FromText: fromText,
		ToText:   toText,
	}
}

func getLineNumberSelector(value []byte) lineNumberSelector {
	fromLine, err := jsonparser.GetInt(value, "fromLine")

	if err != nil {
		fromLine = -1
	}

	toLine, err := jsonparser.GetInt(value, "toLine")

	if err != nil {
		toLine = -1
	}

	return lineNumberSelector{
		FromLine: fromLine,
		ToLine:   toLine,
	}
}

func (rs regexSelector) asContentSelector() (contentSelector, error) {
	compiledRegex, err := regexp.Compile(rs.RegexPattern)

	if err != nil {
		return nil, fmt.Errorf("ERROR: Regex %s for selector did not compile. Error is %s", rs.RegexPattern, err.Error())
	}

	return func(c content) content {
		result := compiledRegex.FindStringSubmatch(c.OriginalText)

		for k, v := range result {
			if int64(k) == rs.GroupNumber {
				newContent := content{OriginalText: v}
				newContent.prepare()
				return newContent
			}
		}

		return content{OriginalText: ""}
	}, nil
}

func (lns lineNumberSelector) asContentSelector() (contentSelector, error) {
	return func(c content) content {
		lines := strings.Split(c.OriginalText, "\n")

		if lns.FromLine == -1 {
			lns.FromLine = 1
		}

		if lns.ToLine > int64(len(lines)) || lns.ToLine == -1 {
			lns.ToLine = int64(len(lines))
		}

		selectedLines := strings.Join(lines[lns.FromLine-1:lns.ToLine], "\n")
		newContent := content{
			OriginalText: selectedLines,
		}
		newContent.prepare()
		return newContent
	}, nil
}

func (tbs textBlockSelector) asContentSelector() (contentSelector, error) {
	return func(c content) content {
		var fromIndex, toIndex int

		if len(c.OriginalText) < 1 {
			return content{OriginalText: ""}
		}

		fromIndex = strings.Index(c.OriginalText, tbs.FromText)
		if fromIndex == -1 || tbs.FromText == "" {
			fromIndex = 0
		}

		toIndex = strings.Index(c.OriginalText, tbs.ToText)
		if toIndex == -1 || tbs.ToText == "" {
			toIndex = len(c.OriginalText) - 1
		}

		if toIndex > len(c.OriginalText) {
			toIndex = len(c.OriginalText) - 1
		}

		newContent := content{OriginalText: c.OriginalText[fromIndex:toIndex]}
		newContent.prepare()

		return newContent
	}, nil
}

func (cs contentSelector) addNestedSelector(wrappingSelector contentSelector) contentSelector {
	return func(c content) content {
		return wrappingSelector(cs(c))
	}
}
