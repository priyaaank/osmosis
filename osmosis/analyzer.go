package osmosis

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

const (
	whiteSpaceRegex  = "\\s+"
	specialCharRegex = "[^a-zA-Z0-9\\s\\.]+"
)

type content struct {
	OriginalText  string
	SanitizedText string
	Words         []string
}

type templates map[string]template

type extractedContent struct {
	AttributeName  string
	AttributeValue string
}
type section struct {
	Selector   contentSelector
	Extractors []contentExtractor
}

type template struct {
	Name     string
	Matcher  contentMatcher
	Sections []section
}

//LoadConfigFile loads the config from a file. The path of the file is provided to the function as the parameter. The config is expected
//in the prescribed json format. It returns a Templates object that contains all template declarations internally.
func LoadConfigFile(filePath string) (templates, error) {
	templateMatchers, _ := ioutil.ReadFile(filePath)
	return LoadConfig(templateMatchers)
}

func LoadConfig(configString []byte) (templates, error) {
	templateErrors := make([]error, 0)
	templates := map[string]template{}
	jsonparser.ArrayEach(configString, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("ERROR: Error extracting template block in the config at index %d. Continuing to next one", offset)
			return
		}

		template, err := ParseTemplate(value)
		if err != nil {
			templateErrors = append(templateErrors, err)
			return
		}

		templates[template.Name] = template
	}, "templates")

	if len(templateErrors) > 0 {
		return nil, templateErrors[0]
	}

	return templates, nil
}

func ParseTemplate(templateDef []byte) (template, error) {
	var templateName string
	var err error
	var newTemplate template
	sectionErrors := make([]error, 0)

	if templateName, err = jsonparser.GetString(templateDef, "templateName"); err != nil {
		return newTemplate, fmt.Errorf("Template name not specified in configuration")
	}

	matcherDef, _, _, err := jsonparser.Get(templateDef, "matchers")

	if err != nil {
		return newTemplate, fmt.Errorf("Matcher block is not specified for template %s. At least one matcher is required for each template", templateName)
	}

	sections := make([]section, 0)
	jsonparser.ArrayEach(templateDef, func(section []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("ERROR: Error extracting sections block in the config at index %d for template %s. Continuing to next one", offset, templateName)
			return
		}

		extractedSection, err := buildSection(section)

		if err != nil {
			sectionErrors = append(sectionErrors, err)
			return
		}

		sections = append(sections, extractedSection)

	}, "sections")

	if len(sectionErrors) > 0 {
		return newTemplate, sectionErrors[0]
	}

	matcher, err := classifyAndBuildMatcher(matcherDef)

	if err != nil {
		return newTemplate, err
	}

	newTemplate.Name = templateName
	newTemplate.Matcher = matcher
	newTemplate.Sections = sections

	return newTemplate, nil
}

func (t *templates) ParseText(docContent string) []extractedContent {

	matchingKeyValues := make([]extractedContent, 0)
	templateMap := map[string]template(*t)
	contentToMatch := content{OriginalText: docContent}
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

func (c *content) prepare() error {
	whtSpaceRegex, err := regexp.Compile(whiteSpaceRegex)

	if err != nil {
		return fmt.Errorf("ERROR: The whitespace replacement regex could not be compiled. Error is %s", err.Error())
	}

	splCharRegex, err := regexp.Compile(specialCharRegex)

	if err != nil {
		return fmt.Errorf("ERROR: The special char replacement regex could not be compiled. Error is %s", err.Error())
	}

	c.SanitizedText = whtSpaceRegex.ReplaceAllString(c.OriginalText, " ")
	c.SanitizedText = splCharRegex.ReplaceAllString(c.SanitizedText, "")
	c.SanitizedText = strings.TrimSpace(c.SanitizedText)
	c.Words = strings.Split(c.SanitizedText, " ")

	return nil
}

func buildSection(value []byte) (section, error) {
	configuredSection := section{}
	selectorSection, _, _, err := jsonparser.Get(value, "contentSelector")
	if err != nil {
		log.Printf("WARN: Could not find configured selector in the section. Continuing assuming extractors will run on full content")
	}

	contentSelector, err := classifyAndBuildSelector(selectorSection)

	if err != nil {
		return configuredSection, fmt.Errorf("ERROR: Could not create selector block for the section. Error is %s", err.Error())
	}

	extractors := make([]contentExtractor, 0)
	var sectionError error

	jsonparser.ArrayEach(value, func(extractor []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			sectionError = err
		}

		parsedExtractor, err := classifyAndBuildExtractor(extractor)

		if err != nil {
			sectionError = err
		}

		extractors = append(extractors, parsedExtractor)
	}, "contentExtractors")

	if sectionError != nil {
		return section{}, sectionError
	}

	return section{
		Selector:   contentSelector,
		Extractors: extractors,
	}, nil
}
