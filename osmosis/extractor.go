package osmosis

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

const (
	whiteSpaceRegex  = "\\s+"
	specialCharRegex = "[^a-zA-Z0-9\\s\\.]+"
)

type Content struct {
	OriginalText  string
	SanitizedText string
	Words         []string
}

type ContentSelector func(c Content) Content

type contentMatcher func(c Content) bool

type Templates map[string]template

type ExtractedContent struct {
	AttributeName  string
	AttributeValue string
}

type contentExtractor func(c Content) ExtractedContent

type section struct {
	Selector   ContentSelector
	Extractors []contentExtractor
}

type template struct {
	Name     string
	Matcher  contentMatcher
	Sections []section
}

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

type regexExtractor struct {
	Regex         string
	AttributeName string
	DefaultValue  string
	GroupNumber   int64
}

func LoadConfigFile(filePath string) Templates {
	templateMatchers, _ := ioutil.ReadFile(filePath)
	return LoadConfig(templateMatchers)
}

func LoadConfig(configString []byte) Templates {
	templates := map[string]template{}
	jsonparser.ArrayEach(configString, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			panic(err)
		}
		template := ParseTemplate(value)
		templates[template.Name] = template
	}, "templates")

	return templates
}

func ParseTemplate(templateDef []byte) template {
	var templateName string
	var err error

	if templateName, err = jsonparser.GetString(templateDef, "templateName"); err != nil {
		panic(err)
	}

	matcherDef, _, _, err := jsonparser.Get(templateDef, "matchers")

	if err != nil {
		panic(err)
	}

	sections := make([]section, 0)
	jsonparser.ArrayEach(templateDef, func(section []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			panic(err)
		}
		sections = append(sections, buildSection(section))
	}, "sections")

	return template{
		Name:     templateName,
		Matcher:  classifyAndBuildMatcher(matcherDef),
		Sections: sections,
	}
}

func (c *Content) prepare() {
	whtSpaceRegex, err := regexp.Compile(whiteSpaceRegex)

	if err != nil {
		panic(err)
	}

	splCharRegex, err := regexp.Compile(specialCharRegex)

	if err != nil {
		panic(err)
	}

	c.SanitizedText = whtSpaceRegex.ReplaceAllString(c.OriginalText, " ")
	c.SanitizedText = splCharRegex.ReplaceAllString(c.SanitizedText, "")
	c.SanitizedText = strings.TrimSpace(c.SanitizedText)
	c.Words = strings.Split(c.SanitizedText, " ")
}

func (t *Templates) ParseText(content string) []ExtractedContent {

	matchingKeyValues := make([]ExtractedContent, 0)
	templateMap := map[string]template(*t)
	contentToMatch := Content{OriginalText: content}
	contentToMatch.prepare()

	for _, template := range templateMap {
		if !template.Matcher(contentToMatch) {
			continue
		}

		for _, section := range template.Sections {
			selectedContent := section.Selector(contentToMatch)

			for _, extractor := range section.Extractors {
				matchingKeyValues = append(matchingKeyValues, extractor(selectedContent))
			}

		}

	}

	return matchingKeyValues
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

func buildSection(value []byte) section {
	selectorSection, _, _, err := jsonparser.Get(value, "contentSelector")
	if err != nil {
		panic(err)
	}
	contentSelector := classifyAndBuildSelector(selectorSection)
	extractors := make([]contentExtractor, 0)

	jsonparser.ArrayEach(value, func(extractor []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			panic(err)
		}
		extractors = append(extractors, classifyAndBuildExtractor(extractor))
	}, "contentExtractors")

	return section{
		Selector:   contentSelector,
		Extractors: extractors,
	}
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

func classifyAndBuildExtractor(value []byte) contentExtractor {
	extractorType, err := jsonparser.GetString(value, "extractorType")

	if err != nil {
		panic(err)
	}

	if strings.EqualFold(extractorType, "regexExtractor") {
		return getRegexExtractor(value).asContentExtractor()
	}

	panic(errors.New("Unknown matcher type"))
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

func (cs ContentSelector) addNestedSelector(wrappingSelector ContentSelector) ContentSelector {
	return func(c Content) Content {
		return wrappingSelector(cs(c))
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
