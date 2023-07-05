[![GoDoc](https://godoc.org/github.com/priyaaank/osmosis/osmosis?status.svg)](https://godoc.org/github.com/priyaaank/osmosis/osmosis)
[![Build Status](https://travis-ci.org/priyaaank/osmosis.svg?branch=master)](https://travis-ci.org/priyaaank/osmosis)
[![Maintainability](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/maintainability)](https://codeclimate.com/github/priyaaank/osmosis/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2ff78eb41e08b7dff42d/test_coverage)](https://codeclimate.com/github/priyaaank/osmosis/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/priyaaank/osmosis)](https://goreportcard.com/report/github.com/priyaaank/osmosis)

# Osmosis

ಜೂ-ಲ್ಯಾಂಗ್ ಲೈಬ್ರರಿ ಜೋಡಿಸಿದ ಹೊಸಕಟ್ಟು ಮತ್ತು ಜೇಡಿಟಿ ಟೆಂಪ್ಲೇಟ್‌ಗಳ ಆಧಾರದ ಮೇಲೆ ಹೊಂದಿಕೊಳ್ಳಲು ಒಂದು ಗ್ರಂಥಾಲಯ.

## Getting started

ಈ ವಿಭಾಗ ನಿಮಗೆ ಒಸ್ಮೋಸಿಸ್ ಫ್ರೇಮ್‌ವರ್ಕ್ ಬಳಸುವುದರ ಪ್ರಾರಂಭದಲ್ಲಿ ನೆರವಾಗುತ್ತದೆ.

### ಪದವನ್ನು ಕನ್ನಡದಲ್ಲಿ "ಸ್ಥಾಪನೆ" ಎಂದು ಅನುವಾದಿಸಲಾಗಿದೆ.

ಒಸ್ಮೋಸಿಸ್ ಅನ್ನು ಸ್ಥಾಪಿಸಲು ನೀವು ಕೆಳಗಿನ ಆದೇಶವನ್ನು ಚಲಾಯಿಸಬಹುದು:

`go get -t github.com/priyaaank/osmosis/osmosis`

### ಬಳಕೆ

ಕಾನ್ಫಿಗ್ ಟೆಂಪ್ಲೇಟ್‌ಗಳನ್ನು ಹೊಸದಾಗಿ ಸ್ಥಳೀಯ ನಿಮಿತ್ತದಲ್ಲಿ ಸಂಗ್ರಹಿಸಲು ನಿಮ್ಮ ಪ್ರೋಗ್ರಾಮದಲ್ಲಿ ಎಲ್ಲಿಯಾದರೂ ಇಟ್ಟಿದ್ದರೆ ಕಾನ್ಫಿಗ್ ಲೋಡ್ ಮಾಡುವಾಗ ಪಥವನ್ನು ನೀಡಲು ಸಾಧ್ಯ. ಕೆಳಗಿನಂತೆ, ಟೆಂಪ್ಲೇಟ್‌ಗಳನ್ನು ಹೊಂದಿರುವ ಕಾನ್ಫಿಗ್ ಉದಾಹರಣೆಯನ್ನು ಹೇಗೆ ಟೆಕ್ಸ್ಟ್ ಫೈಲ್‌ನಿಂದ ವಿಭಜಿಸಿ ಮತ್ತು ಗೆಳೆಯರನ್ನು ಹೊಂದಿಸಬೇಕು ಎಂಬುದರ ಉದಾಹರಣೆ ಕೊಟಿಯಾಗಿದೆ.

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

ವಿಕಲ್ಪವಾಗಿ, ಕಾನ್ಫಿಗ್ ಅನ್ನು `osmosis.LoadConfig()` ಪದ್ಧತಿಗೆ `[]byte` ಆಂಕಣವಾಗಿ ಒಂದು ಇನ್‌ಪುಟ್ ಆಗಿ ಒದಗಿಸಲು ಸಾಧ್ಯ.

### Examples

ನೀವು ಕೆಲವು ಉದಾಹರಣೆಗಳನ್ನು [ಇಲ್ಲಿ](https://github.com/priyaaank/osmosis/tree/master/examples) ಅನುಸರಿಸಿ ಅನುಭವಿಸಬಹುದು.

### Adding a new template

ಕಾನ್ಫಿಗ್ ಫೈಲ್‌ನಲ್ಲಿ ಹೊಸ ಟೆಂಪ್ಲೇಟ್‌ನ್ನು ಸೇರಿಸಲು, ಟೆಂಪ್ಲೇಟ್‌ಗಳ ವಿಭಾಗದಲ್ಲಿ ಹೊಸ ಎಂಟ್ರಿಯನ್ನು ಸೇರಿಸಿ. ಒಂದು ಸರಳ ವಿವರಣೆ ಕೆಳಗಿನಂತೆ ಇರಬಹುದು. ಟೆಂಪ್ಲೇಟ್‌ಗಳನ್ನು ಹೆಚ್ಚುವರಿ ಅಣಿಗಳ ರೂಪದಲ್ಲಿ ಬೇರೆ ಬೇರೆಯಾಗಿ ವಿಭಜಿಸಬೇಕು.

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


ಒಸ್ಮೋಸಿಸ್ ಒಂದು ಗೋ-ಲ್ಯಾಂಗ್ ಪ್ರಕಟಣೆಯ ಲೈಬ್ರರಿ. ಇದು ಜೂಲ್ಯಾಂ ಟೆಂಪ್ಲೇಟ್‌ಗಳ ಆಧಾರದ ಮೇಲೆ ಡೇಟಾವನ್ನು ಜೋಡಿಸುವುದು ಮತ್ತು ಹೊಂದಿಸುವುದರ ಜೊತೆಗೆ ಟೆಕ್ಸ್ಟ್ ಪತ್ರಿಕೆಯಿಂದ ಪಠ್ಯವನ್ನು ಹೊರತೆಗೆಯುವುದಕ್ಕೆ ಜೂಲ್ಯಾಂ ಆಧಾರಿತ ಕಸ್ಟಮ್ ಕಾನ್ಫಿಗ್ಯೂರೇಶನ್ DSL ಅನ್ನು ಬಳಸುತ್ತದೆ. ಒಸ್ಮೋಸಿಸ್ ಟೆಂಪ್ಲೇಟ್‌ಗಳಲ್ಲಿ ಮೂರು ಮುಖ್ಯ ಘಟಕಗಳನ್ನು ಕಾನ್ಫಿಗರ್ ಮಾಡಬೇಕಾಗಿದೆ.
1. ಟೆಂಪ್ಲೇಟ್ ನಾಮ್: ಪ್ರತಿ ಟೆಂಪ್ಲೇಟ್ ಒಂದು ಹೊಸ ಟೆಂಪ್ಲೇಟ್ ಹೆಸರನ್ನು ಹೊಂದಬೇಕು.
2. ಟೆಂಪ್ಲೇಟ್ ಕನಿಷ್ಠ ಅಂಶ: ಟೆಂಪ್ಲೇಟ್ ಮ್ಯಾಚ್ ಮತ್ತು ಡೇಟಾ ಹೊಂದಿಸುವುದಕ್ಕ


* Matcher
* Selector
* Extractor

ಪ್ರತಿ ಒಂದು ಘಟಕವನ್ನು ಕೆಳಗಿನಂತೆ ವಿವರಿಸಲಾಗಿದೆ. ಅದರಲ್ಲಿ ಕೆಲವು ನಮೂದಿಗಳು ಹೊಂದಿವೆ. ಇವುಗಳನ್ನು ಮರೆಮಾಡದೆ ಆಧರಿಸಿ, ಟೆಂಪ್ಲೇಟ್ ಆಧಾರಿತ ಹೊಂದಿಸುವುದನ್ನು ಕಾನ್ಫಿಗರ್ ಮಾಡುವುದು ಸಹಾಯಕವಾಗುತ್ತದೆ.

### Matcher

ಮ್ಯಾಚರ್ ಒಂದು ಕಾನ್ಫಿಗರೇಷನ್ ಬ್ಲಾಕ್ ಆಗಿದೆ ಯಾವುದೋ ಪಠ್ಯದ ಪತ್ರಿಕೆಯನ್ನು ಕಾನ್ಫಿಗರ್ ಮಾಡುತ್ತದೆ. ಒಂದು ಧ್ವನಿಮಾನವನ್ನು ಟೆಂಪ್ಲೇಟ್ ಕಾನ್ಫಿಗರೇಷನ್‌ಗೆ ಅನ್ವಯಿಸುವುದರಿಂದ, ಪಾಠ್ಯಿಕ ಡೇಟಾದಿಂದ ಹೊಂದಿಸುವುದಕ್ಕೆ ಬಳಸುವಂತಹ ಧ್ವನಿಮಾನವು ಪ್ರಕಟವಾಗಿದ್ದರೆ, ಪ್ರಥಮ ಮ್ಯಾಚರ್ ಆಉಟ್‌ಪುಟ್ ನಿಂತಿರುವ ಟೆಂಪ್ಲೇಟ್ ಕಾನ್ಫಿಗರೇಷನ್‌ನಿಂದ ಡೇಟಾವನ್ನು ಹೊಂದಿಸಲು ಬಳಸಲಾಗುತ್ತದೆ. ಪ್ರತಿ ಕಾನ್ಫಿಗರ್ ಮಾಡಲಾದ ಟೆಂಪ್ಲೇಟ್‌ಗೆ ಒಂದು ಮ್ಯಾಚರ್ ಬ್ಲಾಕ್ ಇರುತ್ತದೆ. ಪಠ್ಯದ ಪತ್ರಿಕೆಯನ್ನು ಮ್ಯಾಚರ್ ಬ್ಲಾಕ್‌ಗಳಲ್ಲಿ ಕ್ರಮಿಸಲಾಗುತ್ತ

#### One word matcher

ಇದು ಎಲ್ಲಾ ಮ್ಯಾಚರ್ ಗಳಲ್ಲಿ ಅತ್ಯಲ್ಪ ಪ್ರಕಟವಾದುದು. ಕಾನ್ಫಿಗರೇಷನ್ ನಲ್ಲಿ ಕೊಟ್ಟಿರುವ ಪದಗಳ ಪಟ್ಟಿಗೆ ಅನ್ವಯಿಸಿ, ಒಂದು ಪದವನ್ನು ಹೊಂದಿದ್ದರೆ ಪ್ರತಿಕ್ರಿಯೆಯನ್ನು ಪಡೆಯುತ್ತದೆ. ಪಟ್ಟಿಯಲ್ಲಿರುವ ಪದಗಳು ಕಾಗದ ಪುಸ್ತಕದಲ್ಲಿನ ಪಠ್ಯವನ್ನು ಪ್ರತ್ಯೇಕಿಸುವಲ್ಲಿ, ಪದಗಳನ್ನು ಅಂತರವಿರುವುದುಂಟು ಪ್ರದರ್ಶಿಸುವ ಪಂಕ್ತಿಯಾಗಿ ಹೇಳಬಹುದು. ಒಂದು ಪದದ ಮ್ಯಾಚರ್ ಗಳ ಕಾನ್ಫಿಗರೇಷನ್ ಉದಾಹರಣೆ:

```js
{
    "matcherType": "oneWordMatcher",
    "words": "One upon a time,Harry Potter"
}
```


ಪ್ರಸ್ತುತ ಒಂದು ಪದ ಮ್ಯಾಚರ್ ಅನ್ನು ಅಪ್ರಮೇಯ ಪದ ಪ್ರಮಾಣಿಕತೆ ಮಾಡುವುದಿಲ್ಲ ಮತ್ತು ಪಠ್ಯವನ್ನು ಪೂರ್ಣವಾಗಿ ಪದ ಪದವನ್ನಾಗಿ ಪರೀಕ್ಷಿಸುವುದು ಆವಶ್ಯಕವಾಗಿದೆ. ಉದಾಹರಣೆ 1 ನೇ ಯಾವುದೇ ಒಂದು ಪದ ಮ್ಯಾಚರ್ ಅನ್ನು ಸ್ವೀಕರಿಸುತ್ತದೆ ಹೊಂದಿದೆ, ಆದರೆ ಕೆಳಗಿನ ಉದಾಹರಣೆ 2 ಯಾವುದೇ ಪದ ಮ್ಯಾಚರ್ ಆಗಿಲ್ಲ:

> ಉದಾಹರಣೆ 1: ಪಠ್ಯದಲ್ಲಿ ಎರಡೂ ಅಂಶಗಳು ಪೂರ್ತಿಯಾಗಿ ಕಂಡುಬಂದಿರುವುದರಿಂದ ಮ್ಯಾಚ್ ಆಗುತ್ತದೆ. ಆದರೆ ಎರಡು ಅಂಶಗಳಲ್ಲಿ ಯಾವುದೇ ಒಂದು ಅಂಶವನ್ನು ಮಾತ್ರ ಹೊಂದಿದರೆ ಮ್ಯಾಚ್ ಆಗುತ್ತಿತ್ತು.

```
Once upon a time 
There was a boy called Harry Potter. 
```

> ಉದಾಹರಣೆ 2: ಹ್ಯಾರಿ ಮುಂತಾದ ಪದದ ನಂತರ ಹೊರಡುವ ಹೊಸ ಗೆರೆಯು ಮ್ಯಾಚ್ ಆಗದು. ಇದರ ಜೊತೆಗೆ "Once" ಪದದಲ್ಲಿರುವ "O" ಮತ್ತು "potter" ಪದದಲ್ಲಿರುವ "P" ವರ್ಣಗಳ ಕೀಲಿಕೈಯಿಂದ ಮ್ಯಾಚ್ ಆಗದು.

```
once upon a time
there was a boy called Harry 
potter.
```

