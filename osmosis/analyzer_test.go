package osmosis

import (
	"strings"
	"testing"
)

var contentString = `ANI Technologies Pvt. Ltd. 5th Floor,Infotech Center,, Domlur, Bengaluru, Karnataka 560000 Invoice ID 1IE88NHTQ55547
Customer Name Jacob
Description
Ola Convenience Fee - 1IE88NHTQ55547
Convenience Fee (Ride)
Convenience Fee (Play Convenience Fee(8%))
CGST 9.0%
SGST 9.0%
Total Convenience Fee Fare
Authorised Signatory
`

var testConfig = `{
	"templates": [
		{
			"templateName": "Ola",
			"matchers": {
				"matcherType": "conditionalMatcher",
				"condition": "or",
				"expressions": [
					{
						"matcherType":"conditionalMatcher",
						"condition": "and",
						"expressions": [
							{
								"matcherType": "oneWordMatcher",
								"words": "ANI,Simple,Always"
							},
							{
								"matcherType": "regexMatcher",
								"regexExpression": "ANI\\s+Technologies"
							}
						]
					}
				]
			},
			"sections" : [
                {
                    "contentSelector": {
						"selectorType": "lineNumberSelector",
						"fromLine" : 1,
						"toLine": 2
                    },
                    "contentExtractors": [
                        {
							"extractorType": "regexExtractor",
							"regex": "Invoice ID\s+([A-Z0-9]+)",
							"attributeName": "invoiceNumber",
							"defaultValue":"NA",
							"groupNumber":1
						}
                    ]
                }
            ]
		}
    ]
}
`

func TestThatContentPreparationPopulatesWords(t *testing.T) {
	c := content{OriginalText: contentString}
	c.prepare()

	if len(c.Words) != 41 {
		t.Errorf("Word count was incorrect. Expected %d but got %d", 32, len(c.Words))
	}
}

func TestThatStringIsSanitizedBeforeItIsConvertedToWords(t *testing.T) {
	c := content{OriginalText: contentString}
	c.prepare()

	expected := []string{"ANI", "Technologies", "Pvt.", "Ltd.", "5th", "FloorInfotech", "Center", "Domlur", "Bengaluru", "Karnataka", "560000", "Invoice", "ID", "1IE88NHTQ55547", "Customer", "Name", "Jacob", "Description", "Ola", "Convenience", "Fee", "", "1IE88NHTQ55547", "Convenience", "Fee", "Ride", "Convenience", "Fee", "Play", "Convenience", "Fee8", "CGST", "9.0", "SGST", "9.0", "Total", "Convenience", "Fee", "Fare", "Authorised", "Signatory"}

	if len(expected) != len(c.Words) {
		t.Errorf("Expected %s but got %s words", expected, c.Words)
	}

	for index, element := range c.Words {
		if (strings.Compare(element, expected[index])) != 0 {
			t.Errorf("String was not sanitized in content. Expected %s, but got %s at position %d", expected[index], element, index)
		}
	}
}

func TestThatASanitizedCopyOfContentIsStored(t *testing.T) {
	c := content{OriginalText: contentString}
	c.prepare()

	expected := "ANI Technologies Pvt. Ltd. 5th FloorInfotech Center Domlur Bengaluru Karnataka 560000 Invoice ID 1IE88NHTQ55547 Customer Name Jacob Description Ola Convenience Fee  1IE88NHTQ55547 Convenience Fee Ride Convenience Fee Play Convenience Fee8 CGST 9.0 SGST 9.0 Total Convenience Fee Fare Authorised Signatory"

	if strings.Compare(expected, c.SanitizedText) != 0 {
		t.Errorf("Content not sanitized. Expected %s but got %s", expected, c.SanitizedText)
	}
}

func TestThatTemplateIsBuiltWithMatcherSelectorAndExtractors(t *testing.T) {
	templates, err := LoadConfig(strings.NewReader(testConfig))

	if err != nil {
		t.Errorf("Did not expect error to be returned. But was %s", err.Error())
	}

	if len(templates) != 1 {
		t.Errorf("Expected atleast one template to be returned")
	}

	template := templates["Ola"]

	if template.Name != "Ola" {
		t.Errorf("Expected Ola template to be present")
	}

	if template.Matcher == nil {
		t.Errorf("Expected atleast one matcher to be set in template")
	}

	if template.Sections[0].Selector == nil {
		t.Errorf("Expected atleast one selector to be set in template")
	}

	if len(template.Sections[0].Extractors) != 1 {
		t.Errorf("Expected atleast one extractor to be set in template")
	}
}

func TestThatKeyValueMapIsReturnedForTextWhenMatchingTemplateIsPresent(t *testing.T) {
	templates, err := LoadConfig(strings.NewReader(testConfig))

	if err != nil {
		t.Errorf("Did not expect error to be returned. But was %s", err.Error())
	}

	keyValuePairs, err := templates.ParseText(strings.NewReader(contentString))

	if err != nil {
		t.Errorf("Did not expect error to be raised but was %s", err.Error())
	}

	if len(keyValuePairs) != 1 {
		t.Errorf("Expected at least one key value pair to be returned")
	}

	if keyValuePairs[0].AttributeName != "invoiceNumber" {
		t.Errorf("Expected attribute key to be invoiceNumber")
	}

	if keyValuePairs[0].AttributeValue != "1IE88NHTQ55547" {
		t.Errorf("Expected attribute value to be 1IE88NHTQ55547")
	}

}
