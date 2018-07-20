package osmosis

import (
	"strings"
	"testing"
)

func TestThatRegexExtractorIsUsedToExtractTheValue(t *testing.T) {
	expectedAttributeName := "invoiceNumber"
	expectedAttributeValue := "1IE88NHTQ55547"
	invoiceExtractor := `{
		"extractorType": "regexExtractor",
		"regex": "Invoice ID\s+([A-Z0-9]+)",
		"attributeName": "invoiceNumber",
		"defaultValue":"NA",
		"groupNumber":1
	}`
	c := Content{OriginalText: "Bengaluru, Karnataka 560000 Invoice ID 1IE88NHTQ55547"}
	c.prepare()

	extractor := classifyAndBuildExtractor([]byte(invoiceExtractor))

	extractedContent := extractor(c)

	if strings.Compare(extractedContent.AttributeName, expectedAttributeName) != 0 {
		t.Errorf("Expected attribute name [%s] to match [%s]", extractedContent.AttributeName, expectedAttributeName)
	}

	if strings.Compare(extractedContent.AttributeValue, expectedAttributeValue) != 0 {
		t.Errorf("Expected attribute name [%s] to match [%s]", extractedContent.AttributeValue, expectedAttributeValue)
	}
}
