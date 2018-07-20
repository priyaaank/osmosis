package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/priyaaank/osmosis/osmosis"
)

func main() {
	inputFilePath, err := filepath.Abs("textfiles/freshmenu_receipt.txt")

	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadFile(inputFilePath)

	if err != nil {
		panic(err)
	}

	confFilePath, err := filepath.Abs("conf/contentMatchers.json")

	if err != nil {
		panic(err)
	}

	templates := osmosis.LoadConfigFile(confFilePath)

	extractedInfo := templates.ParseText(string(fileContent))

	for _, info := range extractedInfo {
		fmt.Printf("AttrName: %s | AttrValue: %s", info.AttributeName, info.AttributeValue)
		fmt.Println("")
	}
}
