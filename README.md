[![Build Status](https://travis-ci.org/priyaaank/osmosis.svg?branch=master)](https://travis-ci.org/priyaaank/osmosis)
[![Maintainability](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/maintainability)](https://codeclimate.com/github/priyaaank/osmosis/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/test_coverage)](https://codeclimate.com/github/priyaaank/osmosis/test_coverage)

# Osmosis

A go-lang library to match and extract data based on json templates.

## Overview

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

```text
Once upon a time 
There was a boy called Harry Potter. 
```

> Example 2 : Does not match because of new line after harry and different case of "O" in Once, and "P" in potter

```text
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

```text
Once upon a time 
There was a boy called Harry Potter. 
```

> Example 2 : Does not match because of the case difference of "H" in harry.

```text
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

```
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



## Getting started


### Installation


### Adding a new template


### Usage


### Examples


### Reference
