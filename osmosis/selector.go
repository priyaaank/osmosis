package osmosis

import (
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

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

func classifyAndBuildSelector(value []byte) ContentSelector {
	var contentSelector ContentSelector

	selectorType, err := jsonparser.GetString(value, "selectorType")

	if err != nil {
		return nil
	}

	if strings.EqualFold(selectorType, "textBlockSelector") {
		contentSelector = getTextBlockSelector(value).asContentSelector()
	} else if strings.EqualFold(selectorType, "lineNumberSelector") {
		contentSelector = getLineNumberSelector(value).asContentSelector()
	} else if strings.EqualFold(selectorType, "regexSelector") {
		contentSelector = getRegexSelector(value).asContentSelector()
	}

	contentSelectorValue, _, _, err := jsonparser.Get(value, "contentSelector")

	if err != nil {
		return contentSelector
	}

	if nestedSelector := classifyAndBuildSelector(contentSelectorValue); nestedSelector != nil {
		contentSelector = contentSelector.addNestedSelector(nestedSelector)
	}

	return contentSelector
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

func (rs regexSelector) asContentSelector() ContentSelector {
	return func(c Content) Content {
		reg, err := regexp.Compile(rs.RegexPattern)

		if err != nil {
			panic(err)
		}

		result := reg.FindStringSubmatch(c.OriginalText)

		for k, v := range result {
			if int64(k) == rs.GroupNumber {
				newContent := Content{OriginalText: v}
				newContent.prepare()
				return newContent
			}
		}

		return Content{OriginalText: ""}
	}
}

func (lns lineNumberSelector) asContentSelector() ContentSelector {
	return func(c Content) Content {
		lines := strings.Split(c.OriginalText, "\n")

		if lns.FromLine == -1 {
			lns.FromLine = 1
		}

		if lns.ToLine > int64(len(lines)) || lns.ToLine == -1 {
			lns.ToLine = int64(len(lines))
		}

		selectedLines := strings.Join(lines[lns.FromLine-1:lns.ToLine], "\n")
		newContent := Content{
			OriginalText: selectedLines,
		}
		newContent.prepare()
		return newContent
	}
}

func (tbs textBlockSelector) asContentSelector() ContentSelector {
	return func(c Content) Content {
		var fromIndex, toIndex int

		if len(c.OriginalText) < 1 {
			return Content{OriginalText: ""}
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

		newContent := Content{OriginalText: c.OriginalText[fromIndex:toIndex]}
		newContent.prepare()

		return newContent
	}
}

func (cs ContentSelector) addNestedSelector(wrappingSelector ContentSelector) ContentSelector {
	return func(c Content) Content {
		return wrappingSelector(cs(c))
	}
}
