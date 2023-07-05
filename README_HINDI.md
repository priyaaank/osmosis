[![GoDoc](https://godoc.org/github.com/priyaaank/osmosis/osmosis?status.svg)](https://godoc.org/github.com/priyaaank/osmosis/osmosis)
[![Build Status](https://travis-ci.org/priyaaank/osmosis.svg?branch=master)](https://travis-ci.org/priyaaank/osmosis)
[![Maintainability](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/maintainability)](https://codeclimate.com/github/priyaaank/osmosis/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/test_coverage)](https://codeclimate.com/github/priyaaank/osmosis/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/priyaaank/osmosis)](https://goreportcard.com/report/github.com/priyaaank/osmosis)

# Osmosis

जस्टेम्पलेट आधारित JSON पर मिलान और डेटा निकालने के लिए एक गो-लैंग पुस्तकालय। (Justemplate नामक गो-लैंग पुस्तकालय)

## Getting started

यह खंड आपको Osmosis फ़्रेमवर्क की शुरुआत में मदद करेगा।

### Installation

Osmosis को स्थापित करने के लिए आप निम्नलिखित कमांड चला सकते हैं।

`go get -t github.com/priyaaank/osmosis/osmosis`

### Usage

कॉन्फ़िग टेम्पलेट किसी भी स्थान पर आपके प्रोग्राम में संग्रहीत किए जा सकते हैं, जब तक फ़ाइल का पथ प्राविधिकता के दौरान कॉन्फ़िगरेशन लोड करने के लिए प्रदान किया जाता हो। निम्नलिखित एक त्वरित उदाहरण है कि कैसे टेम्पलेट्स को संकलित करके एक टेक्स्ट फ़ाइल से फ़ील्ड पार्स और निकाला जा सकता है।

```go
package main

func main() {

    confFile, _ := os.Open("/some/path/on/disk/project/config/osmosisconfig.json")
    contentFile, _ := os.Open("/some/path/on/disk/project/inputfiles/sample.txt")
    templates, err := osmosis.LoadConfig(bufio.NewReader(confFile))
    extractedInfo, err := templates.ParseText(bufio.NewReader(contentFile))
    
    for _, info := range extractedInfo {
        fmt.Printf("AttrName: %s | AttrValue: %s \n", info.AttributeName, info.AttributeValue)
    }

}
```

वैकल्पिक रूप से, कॉन्फ़िग को osmosis.LoadConfig() विधि में []byte इनपुट के रूप में भी प्रदान किया जा सकता है।

### Examples

आप [यहां](https://github.com/priyaaank/osmosis/tree/master/examples) कई उदाहरण देख सकते हैं जो अमल में लाए गए हैं।

### Adding a new template

कॉन्फ़िग फ़ाइल में एक नया टेम्पलेट जोड़ने के लिए, टेम्पलेट सेक्शन में एक नया एंट्री जोड़ें। एक सरल परिभाषा निम्नलिखित रूप में दिखेगी। टेम्पलेट्स को कॉमा से अलग किया जाना चाहिए जैसे कि एक एरे के कई तत्वों के रूप में।

```js
{
    "templateName": "FreshMenu",
    "matchers": {
        "matcherType": "oneWordMatcher",
        "words": "Serendipity,Shanghai"
        }
    },
    "sections" : [
        {
            "contentSelector": {
                "selectorType": "textBlockSelector",
                "fromText" : "CUSTOMER DETAILS",
                "toText": "HSN Code"
            },
            "contentExtractors": [
                {
                    "extractorType": "regexExtractor",
                    "regex": "Name:\s+([A-z\s]+)\n",
                    "attributeName": "name",
                    "defaultValue": "NA",
                    "groupNumber": 1
                }
            ]
        }
    ]
}
```

## Overview & examples

Osmosis एक गो-लैंग लाइब्रेरी है जो JSON टेम्पलेट के आधार पर डेटा मिलान और निकालने के लिए लिखी गई है। यह एक टेक्स्ट दस्तावेज़ से पाठ को मिलान और निकालने के लिए JSON आधारित कस्टम कॉन्फ़िगरेशन DSL का उपयोग करता है। Osmosis में प्रत्येक टेम्पलेट में तीन मुख्य घटक होते हैं जिन्हें कॉन्फ़िगर किया जाना चाहिए।

* मिलानकर्ता (Matcher)
* चयनकर्ता (Selector)
* निकालक (Extractor)

निम्नलिखित प्रत्येक को एक संग्रहीत करके नीचे समझाया गया है। इन्हें बेहतर समझने से आप टेम्पलेट आधारित निकालने को कॉन्फ़िगर कर सकते हैं।

### मिलानकर्ता (Matcher)

मिलानकर्ता (Matcher) एक कॉन्फ़िगरेशन ब्लॉक है जो एक पाठिक दस्तावेज़ को कॉन्फ़िगर किए गए टेम्पलेट के साथ जोड़ता है। सकारात्मक मिलान कॉन्फ़िगरेशन पाठ के लिए पाठ से डेटा निकालने के लिए लागू होता है। प्रत्येक कॉन्फ़िगर किए गए टेम्पलेट में एक मिलानकर्ता ब्लॉक होता है। एक पाठिक दस्तावेज़ को सभी मिलानकर्ता ब्लॉकों से क्रमशः पास किया जाता है। पहले सकारात्मक मिलानकर्ता के आउटपुट से डेटा निकालने के लिए निकालने की कॉन्फ़िगरेशन का उपयोग किया जाता है। मिलानकर्ता कई प्रकार के हो सकते हैं।

#### एक शब्द मिलानकर्ता (One word matcher)

यह सभी मिलानकर्ताओं में से सबसे सरल मिलानकर्ता है। कॉन्फ़िगरेशन में दिए गए शब्दों की सूची के लिए, यह उनमें से किसी एक शब्द के साथ मिलान करता है। यदि प्रदान की गई इनपुट दस्तावेज़ में इनमें से एक शब्द भी मिल जाता है, तो यह सकारात्मक मिलान लौटाता है। सूची में शब्दों को अल्पविराम (,) द्वारा अलग किए गए स्थान युक्त पाठ टुकड़ों के रूप में दिया जा सकता है। एक शब्द मिलानकर्ता के लिए एक कॉन्फ़िगरेशन का उदाहरण निम्नलिखित है:

```js
{
    "matcherType": "oneWordMatcher",
    "words": "One upon a time,Harry Potter"
}
```

वर्तमान में एक शब्द मिलानकर्ता अवधिक वर्णमाला संवेदनशील मिलान नहीं करता है। इसके अलावा, इसे पूरी तरह से शब्दों के अनुसार मिलान करने की आवश्यकता होती है। उदाहरण के रूप में नीचे दिए गए दो उदाहरण में, उदाहरण 1 सकारात्मक मिलान है जबकि उदाहरण 2 नहीं है:

> Example 1: Matches as both fragments are found verbatim in the text below. It would have been enough to match only either of the two fragments though.

```
Once upon a time 
There was a boy called Harry Potter. 
```

> Example 2 : Does not match because of new line after harry and different case of "O" in Once, and "P" in potter

```
once upon a time
there was a boy called Harry 
potter.
```

#### All words matcher

Similar to one word matcher, all words matcher requires that all words or text fragments be found in the provided text for a positive match. If all of the words are found in the provided input document, it returns a positive match. The words in the list can be space containing text fragments, seperated by a comma (`,`). An example of configuration for a All words matcher is as follows:

```js
{
    "matcherType": "allWordsMatcher",
    "words": "One upon a time,Harry Potter"
}
```

Currently one word matcher **does not do** case insensitive match. Also, it needs the text to match in it entirity word by word. for instance the example 1 is a positive match where as example 2 below is not. 

> Example 1: Matches as both fragments are found verbatim in the text below.

```
Once upon a time 
There was a boy called Harry Potter. 
```

> Example 2 : Does not match because of the case difference of "H" in harry.

```
Once upon a time
there was a boy called harry Potter.
```

#### Regex matcher

Similar to previous matcher, regex matcher, matches text in the provided input. It takes a single regex expression at a time and returns a positive match indicator if the regex finds a match. 

Sample configuration for regex matcher looks like as follows:

```js
{
    "matcherType": "regexMatcher",
    "words": "(H|h)arry\s[A-z]{6}"
}
```

The above sample configuration will match text containing both `Harry Potter` and `harry potter` but not `Harry P0tter`

#### Conditional matcher

Conditional matcher block is analogous to a logical programatic condition. It supports onlu `AND` and `OR` operations. Using Conditional matcher all other matchers can be grouped together to form sophisticated conditional logic within matcher configuration of a template. Here is a simple example below around how a matcher configuration for following condition can be written:

> Conditions to evaluate on a textual document

* The document contains any of the words [`Serendipity`,`Shanghai`] or contains all words [`29BBZZF8899Q0ZQ`, `U15209KA2014PTC075887`]
* In addition to the first condition, it should also contain at least one of the text fragment `HSR Layout`,`orders@freshmenu.com`

If both conditions are true above, then it matches the template name `Freshmenu`

> Sample configuration in JSON DSL

```js
"matchers": {
    "matcherType": "conditionalMatcher",
    "condition": "and",
    "expressions": [
        {
            "matcherType":"conditionalMatcher",
            "condition": "or",
            "expressions": [
                {
                    "matcherType": "oneWordMatcher",
                    "words": "Serendipity,Shanghai"
                },
                {
                    "matcherType": "allWordsMatcher",
                    "words": "29BBZZF8899Q0ZQ,U15209KA2014PTC075887"
                }
            ]
        },
        {
            "matcherType": "oneWordMatcher",
            "words": "HSR Layout,orders@freshmenu.com"
        }
    ]
}
```

### Sections

Once input text has been matched to a configured template using a selector config, the sections of the template are used to select and extract the text. Sections is a list of sections. Each section can contain a single `Selector` and a list of `Extractors`. Each `selector` block selects the part of provided text input. The selector text block is handed over to the extractors. Each extractor extracts the targetted text and returns it in a key value pair format. Both `Selector` and `Extractor` are explained in more details in sections below.

### चयनकर्ता (Selector)

A selector is a configuration block in JSON DSL config which selects a part of the content from the provided input. There will be only one block of selector in a Section, however the selectors can be nested. Each selector obtains the selected block of text from the previous selector and operates on it to provide a selected block of text. There are several types of selectors as explained below.

Each Selector can have a optional child config element called `contentSelector`. This contains a nested selector that operates on the content provided by previous selector. This is how selectors can be nested to do multiple levels of selection.

#### TextBlockSelector

Text block selector is a simple selector that takes `fromText` and `toText` attribute. It finds the first occurance of `fromText` and selects text including the text of fromText. It finds the first occurance of `toText` and selects all the text in between. This selection excludes the text specified in `toText`. If the `fromText` tag is missing from the definition, then content from the beginning is selected. If the `toText` tag is missing, content is selected till be end of the document by default. 

> A sample configuration for TextBlockSelector will be:

```js
"contentSelector": {
    "selectorType": "textBlockSelector",
    "fromText" : "Listening",
    "toText": "Issa"
}
```

> Example content

```
Winter seclusion -
Listening, that evening,
To the rain in the mountain.
- Kobayashi Issa
```

> Output of selector

```
Listening, that evening,
To the rain in the mountain.
- Kobayashi
```

#### LineNumberSelector

This selects the content from the `fromLine`, including the content of `fromLine` till the content of `toLine`, including the content of `toLine`. If the start line is in negative somehow, the content from the first line is selected. Similarly, if the `toLine` is not in range, the content till last line is selected.  

There are few things that should be noted around LineNumberSelector

* When nested, uses the line number for the content selected by previous selector and not original selector.
* The line numbers should be integers in config. If they are not, currently program does not raise an error but processes them by substituting the line numbers to be 0. A fix will be made to address this issue soon. 

> Example : Sample content below

```
Winter seclusion -
Listening, that evening,
To the rain in the mountain.
- Kobayashi Issa
```

> Configured nested line selector

```js
"contentSelector": {
    "selectorType": "textBlockSelector",
    "fromText" : "Listening",
    "toText": "Issa",
    "contentSelector" : {
        "selectorType": "lineNumberSelector",
        "fromLine": 1,
        "toLine": 2
    }
}
```

In the above example the first TextBlockSelector will select following lines.

```
Listening, that evening,
To the rain in the mountain.
- Kobayashi 
```

The LineBlockSelector will treat line 1 as `Listening, that evening,` and line 2 as `To the rain in the mountain.`. So original line numbers are not applicable with a nested lineNumberSelector. 

> Final output

```
Listening, that evening,
To the rain in the mountain.
```

#### RegexSelector

Regex Selector is a selector that uses a regex expression to select a block of text. It can be nested as needed with other selectors. A regex selector also takes a group number along with the regex expression. If there are multiple groups of matches, the selector can specify which group should be selected. 

A sample configuration for RegexSelector looks as follows

```js
"contentSelector": {
    "selectorType": "regexSelector",
    "regex" : "Listening,([\w\s]+),",
    "groupNumber": 1,
}
```

> Input text

```
Winter seclusion -
Listening, that evening,
To the rain in the mountain.
- Kobayashi Issa
```

> Output from the previous selector

```
that evening
```

### निकालक (Extractor)

An extractor is a config block that extracts the content for a given key value. The extractor is part of the section in config. The text upon which the extract operates is the one selected by the selector block of the section. This is the final stage of the template extraction. Each extractor works to extract a value for a given key in the cofig. There is an option to provide a default value as well in case a matching value cannot be extracted. Currently extractors do not raise errors when a match is not found in the content and quietly substitute the default value. 

As per the design of JSON config DSL multiple extractors can act and extract data from a selector within each section. All the section contribute to the same key map. The output key map store is flat and contains the keys with their corresponding values extracted by extractors.

#### Regex extractor

Currently the framework supports on a regex extractor. It takes a regex pattern and a group number along with a default value and key name. Group number indicates which matching group should be selected to populate the value.

A sample extractor config looks like as follows. 

```js
{
    "extractorType": "regexExtractor",
    "regex": "(FM-KA-[\d]+)",
    "attributeName": "invoiceNumber",
    "defaultValue": "NA",
    "groupNumber": 1
},
```

> Sample input content. This will be usually an output of a selector config block

```
Invoice No FM-KA-4931389 generated on 12/01/2018
```

> The output from the template selection will be a key value pair of 

```
{
    "invoiceNumber": "FM-KA-4931389"
}
```
### Complete sample config

This is how a sample config looks like with all elements in place.

```js
{
    "templates": [
        {
            "templateName": "FreshMenu",
            "matchers": {
                "matcherType": "conditionalMatcher",
                "condition": "and",
                "expressions": [
                    {
                        "matcherType":"conditionalMatcher",
                        "condition": "or",
                        "expressions": [
                            {
                                "matcherType": "oneWordMatcher",
                                "words": "Serendipity,Shanghai"
                            },
                            {
                                "matcherType": "allWordsMatcher",
                                "words": "29BBZZF8899Q0ZQ,U15209KA2014PTC075887"
                            }
                        ]
                    },
                    {
                        "matcherType": "oneWordMatcher",
                        "words": "HSR Layout,orders@freshmenu.com"
                    }
                ]
            },
            "sections" : [
                {
                    "contentSelector": {
                        "selectorType": "textBlockSelector",
                        "fromText" : "CUSTOMER DETAILS",
                        "toText": "HSN Code",
                        "contentSelector" : {
                            "selectorType": "lineNumberSelector",
                            "fromLine": 1,
                            "toLine": 14,
                            "contentSelector": {
                                "selectorType":"regexSelector",
                                "regex": "[\w\W]+",
                                "groupNumber": 0
                            }
                        }
                    },
                    "contentExtractors": [
                        {
                            "extractorType": "regexExtractor",
                            "regex": "Name:\s+([A-z\s]+)\n",
                            "attributeName": "name",
                            "defaultValue": "NA",
                            "groupNumber": 1
                        },
                        {
                            "extractorType": "regexExtractor",
                            "regex": "(FM[\d]+)",
                            "attributeName": "invoiceNumber",
                            "defaultValue": "NA",
                            "groupNumber": 1
                        },
                        {
                            "extractorType": "regexExtractor",
                            "regex": "\n([\d]+)\n",
                            "attributeName": "phoneNumber",
                            "defaultValue": "NA",
                            "groupNumber": 1
                        }
                    ]
                }
            ]
        },
        {
            "templateName": "UberIndia",
            "matchers": {
                "matcherType": "oneWordMatcher",
                "words": "Uber India Systems,Invoice issued by Uber"
            },
            "sections" : [
                {
                    "contentSelector": {
                        "selectorType": "textBlockSelector",
                        "fromText" : "Invoice Number",
                        "toText": "Tax Amount",
                        "contentSelector" : {
                            "selectorType": "lineNumberSelector",
                            "fromLine": 1,
                            "toLine": 10
                        }
                    },
                    "contentExtractors": [
                        {
                            "extractorType": "regexExtractor",
                            "regex": "Invoice\s+Number:\s+([a-zA-Z0-9]+-[0-9]+-[0-9]+-[0-9]+)",
                            "attributeName": "invoiceNumber",
                            "defaultValue": "NA",
                            "groupNumber": 1
                        },
                        {
                            "extractorType": "regexExtractor",
                            "regex": "Invoice issued by Uber[\S\s]+:\n([a-zA-Z\s]+)\n",
                            "attributeName": "driverName",
                            "defaultValue": "NA",
                            "groupNumber": 1
                        }
                    ]
                },
                {
                    "contentSelector": {
                        "selectorType": "textBlockSelector",
                        "fromText" : "Gross Amount",
                        "toText": "Category of services",
                        "contentSelector" : {
                            "selectorType": "lineNumberSelector",
                            "fromLine": 4
                        }
                    },
                    "contentExtractors": [
                        {
                            "extractorType": "regexExtractor",
                            "regex": "(\d+.\d+)[\D\W\s]+(\d+.\d+)[\D\W\s]+(\d+.\d+)",
                            "attributeName": "totalAmount",
                            "defaultValue": "NA",
                            "groupNumber": 3
                        }
                    ]
                }
            ]
        }
    ]
}

```

## Developer setup for contribution

### Installing

Clone the repo

`git clone git@github.com:priyaaank/osmosis.git`

Install glide package manager

`go get -t github.com/Masterminds/glide`

Install dependencies

`glide install`

### Running example

Change to examples dir

`cd examples`

Run the parser file

`go run parser.go`
