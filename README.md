[![Build Status](https://travis-ci.org/priyaaank/osmosis.svg?branch=master)](https://travis-ci.org/priyaaank/osmosis)
[![Maintainability](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/maintainability)](https://codeclimate.com/github/priyaaank/osmosis/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/test_coverage)](https://codeclimate.com/github/priyaaank/osmosis/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/priyaaank/osmosis)](https://goreportcard.com/report/github.com/priyaaank/osmosis)

# Osmosis

A go-lang library to match and extract data based on json templates.

## Getting started

This section will help you get started with Osmosis framework. 

### Installation

To install osmosis you can run following command

`go get -t github.com/priyaaank/osmosis/osmosis`

### Usage

Config templates can be stored anywhere in your program as long as the path to file is provided to load the configuration while initialization. Following is a quick example of how a config containing templates can be used to parse and extract fields from a text file. 

```go
package main

func main() {

    templates := osmosis.LoadConfigFile("/some/path/on/disk/project/config/osmosisconfig.json")
    contentToParse, _ :=  ioutil.ReadFile("/some/path/on/disk/project/inputfiles/sample.txt")
    extractedContent := templates.ParseText(string(contentToParse))
    
    for _, info := range extractedInfo {
        fmt.Printf("AttrName: %s | AttrValue: %s \n", info.AttributeName, info.AttributeValue)
    }

}
```

Alternatively the config can also be provided as an `[]byte` input to the `osmosis.LoadConfig()` method.  

### Examples

You can find several examples implemented [here](https://github.com/priyaaank/osmosis/tree/master/examples)

### Adding a new template

To add a new template in the config file, add a new enrty in templates section. A simple definition would look like as follows. The templates should be separated by commas as multiple elements of an array.

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

Osmosis is a library written in go-lang to match and extract data based on json templates. It uses a JSON based custom configuration DSL to build templates that can match and extract text from a textual document. Osmosis has three key components in each template that need to be configured. 

* Matcher
* Selector
* Extractor

Each of them is explained below with few sample configurations. Understanding them better will help you configure a template based extraction. 

### Matcher

Matcher is a block of configuration that pairs a textual document with a configured template. A positive match applies the template configuration to the textual data for extraction. Each configured template has a matcher block. A textual document is passed through all the matcher blocks sequentially. The extraction configuration from the first positive matcher output is used to extract data. Matchers can be of several types. 

#### One word matcher

This is one of the simplest matcher of all. For a list of given words in configuration, it matches one of the words. If even one of the words is found in the provided input document, it returns a positive match. The words in the list can be space containing text fragments, seperated by a comma (`,`). An example of configuration for a One word matcher is as follows:

```js
{
    "matcherType": "oneWordMatcher",
    "words": "One upon a time,Harry Potter"
}
```

Currently one word matcher does not do case insensitive match. Also, it needs the text to match in it entirity word by word. for instance the example 1 is a positive match where as example 2 below is not. 

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

### Selectors

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

### Extractors

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