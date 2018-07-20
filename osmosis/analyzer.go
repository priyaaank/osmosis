package osmosis

import (
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

type Templates map[string]template

type ExtractedContent struct {
	AttributeName  string
	AttributeValue string
}
type section struct {
	Selector   ContentSelector
	Extractors []contentExtractor
}

type template struct {
	Name     string
	Matcher  contentMatcher
	Sections []section
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
