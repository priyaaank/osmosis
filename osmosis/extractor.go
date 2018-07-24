package osmosis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type contentExtractor func(c content) extractedContent

type regexExtractor struct {
	Regex         string
	AttributeName string
	DefaultValue  string
	GroupNumber   int64
}

func classifyAndBuildExtractor(value []byte) (contentExtractor, error) {
	extractorType, err := jsonparser.GetString(value, "extractorType")

	if err != nil {
		return nil, err
	}

	if strings.EqualFold(extractorType, "regexExtractor") {
		return getRegexExtractor(value).asContentExtractor()
	}

	return nil, fmt.Errorf("ERROR: Unknown extractor type %s", extractorType)
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

func (re regexExtractor) asContentExtractor() (contentExtractor, error) {
	compiledRegex, err := regexp.Compile(re.Regex)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Could not compile the extractor regex %s", re.Regex)
	}

	return func(c content) extractedContent {

		extractedKeyVal := extractedContent{
			AttributeName:  re.AttributeName,
			AttributeValue: re.DefaultValue,
		}

		result := compiledRegex.FindStringSubmatch(c.OriginalText)

		for k, val := range result {
			if int64(k) == re.GroupNumber {
				extractedKeyVal.AttributeValue = strings.TrimSpace(val)
				return extractedKeyVal
			}
		}

		return extractedKeyVal
	}, nil
}
