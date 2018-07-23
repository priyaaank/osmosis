package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/priyaaank/osmosis/osmosis"
)

func main() {
	fmt.Println("********************  Uber Receipt *****************************")
	extractData("textfiles/uber_india_receipt.txt")
	fmt.Println("******************** Uber Receipt  *****************************\n")

	fmt.Println("******************** Freshmenu Receipt *****************************")
	extractData("textfiles/freshmenu_receipt.txt")
	fmt.Println("******************** Freshmenu Receipt *****************************\n")
}

func extractData(inputFilePath string) {
	absFilePath, err := filepath.Abs(inputFilePath)

	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadFile(absFilePath)

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
		fmt.Printf("AttrName: %s | AttrValue: %s \n", info.AttributeName, info.AttributeValue)
	}
}
