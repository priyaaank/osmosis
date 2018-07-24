//Package osmosis allows configured templates to match, select and extract attribute values from textual files
package osmosis

import (
	"fmt"
	"io"
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

//Templates is internally a map of configured templates. The key is the name of the template and the value is a template struct object.
//Method that utilize configured templates can be called on this struct.
type Templates map[string]template

//ExtractedContent is an object which represents a key value pair. For each configured extractors an ExtractedContent can be returned.
//AttributeName represents the configured key for the pair.
//AttributeValue represents the extracted value for the pair.
type ExtractedContent struct {
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

//LoadConfig loads the configuration from the provided io.Reader object. It expects the content to be in JSON DSL format as explained in docs.
//Once loaded, it creates an internal struct containing all relevant information and returns a Templates object.
//Templates object represent a set of configured templates. Method on this object can be called to parse content to match, select and extract.
//An error can also be returned when config parsing encounters a problem either with minimum required configuration, syntax invalidity or other errors.
func LoadConfig(reader io.Reader) (Templates, error) {
	configString, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	templateErrors := make([]error, 0)
	templates := map[string]template{}
	jsonparser.ArrayEach(configString, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("ERROR: Error extracting template block in the config at index %d. Continuing to next one", offset)
			return
		}

		template, err := parseTemplate(value)
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

//ParseText takes in a io.Reader object that can provide the content that needs to be matched across templates and then extracted from.
//It sequentially runs matchers from all templates configured in system. Once a template matches, it applies the selectors and extractors
//to extract the key value pairs. These key-value pairs are returned as a slice of ExtractedContent.
//[]ExtractedContent represents a slice of all key-value pairs
//This method can also return error if there is a problem while parsing the content with the matched template or when a matching template
//is not found.
func (t *Templates) ParseText(docReader io.Reader) ([]ExtractedContent, error) {

	docContent, err := ioutil.ReadAll(docReader)

	if err != nil {
		return nil, err
	}

	matchingKeyValues := make([]ExtractedContent, 0)
	templateMap := map[string]template(*t)
	contentToMatch := content{OriginalText: string(docContent)}
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

	return matchingKeyValues, nil
}

func parseTemplate(templateDef []byte) (template, error) {
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
