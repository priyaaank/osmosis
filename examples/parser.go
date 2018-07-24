package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/priyaaank/osmosis/osmosis"
)

func main() {
	fmt.Println("********************  Uber Receipt *****************************")
	extractData("textfiles/uber_india_receipt.txt")
	fmt.Println("******************** Uber Receipt  *****************************")
	fmt.Println("")
	fmt.Println("******************** Freshmenu Receipt *****************************")
	extractData("textfiles/freshmenu_receipt.txt")
	fmt.Println("******************** Freshmenu Receipt *****************************")
}

func extractData(inputFilePath string) {
	absFilePath, err := filepath.Abs(inputFilePath)

	if err != nil {
		panic(err)
	}

	contentFile, err := os.Open(absFilePath)

	if err != nil {
		panic(err)
	}

	confFilePath, err := filepath.Abs("conf/contentMatchers.json")

	if err != nil {
		panic(err)
	}

	file, err := os.Open(confFilePath)

	if err != nil {
		panic(err)
	}

	templates, err := osmosis.LoadConfig(bufio.NewReader(file))

	if err != nil {
		panic(err)
	}

	extractedInfo, err := templates.ParseText(bufio.NewReader(contentFile))

	if err != nil {
		panic(err)
	}

	for _, info := range extractedInfo {
		fmt.Printf("AttrName: %s | AttrValue: %s \n", info.AttributeName, info.AttributeValue)
	}
}
