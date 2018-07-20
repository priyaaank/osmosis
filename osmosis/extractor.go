package osmosis

import (
	"errors"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type contentExtractor func(c Content) ExtractedContent

type regexExtractor struct {
	Regex         string
	AttributeName string
	DefaultValue  string
	GroupNumber   int64
}

func classifyAndBuildExtractor(value []byte) contentExtractor {
	extractorType, err := jsonparser.GetString(value, "extractorType")

	if err != nil {
		panic(err)
	}

	if strings.EqualFold(extractorType, "regexExtractor") {
		return getRegexExtractor(value).asContentExtractor()
	}

	panic(errors.New("Unknown extractor type"))
}

func getRegexExtractor(value []byte) regexExtractor {
	regex, _, _, _ := jsonparser.Get(value, "regex")
	attributeName, _ := jsonparser.GetString(value, "attributeName")
	defaultValue, _ := jsonparser.GetString(value, "defaultValue")
	groupNumber, _ := jsonparser.GetInt(value, "groupNumber")

	return regexExtractor{
		Regex:         string(regex),
		AttributeName: attributeName,
		DefaultValue:  defaultValue,
		GroupNumber:   groupNumber,
	}
}

func (re regexExtractor) asContentExtractor() contentExtractor {
	return func(c Content) ExtractedContent {
		reg, err := regexp.Compile(re.Regex)

		extractedContent := ExtractedContent{
			AttributeName:  re.AttributeName,
			AttributeValue: re.DefaultValue,
		}

		if err != nil {
			panic(err)
		}

		result := reg.FindStringSubmatch(c.OriginalText)

		for k, val := range result {
			if int64(k) == re.GroupNumber {
				extractedContent.AttributeValue = val
				return extractedContent
			}
		}

		return extractedContent
	}
}
