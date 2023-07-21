package main

import (
	"log"
	"os"
)

const ERROR = "ERROR: "

var logger = log.Logger{}

func InitLogging() {
	logger.SetFlags(log.Ldate | log.Ltime)
	logger.SetOutput(os.Stdout)
}

func main() {
	InitLogging()

	const link = "https://www.sephora.ro/"

	extractor := HtmlExtractor{link: link}
	hp, err := extractor.GetHtmlPage()
	if err != nil {
		logger.Fatalf("Could not get the html page requested")
	}

	node, err := ParseHtml(hp)
	if err != nil {
		logger.Fatalf("Could not parse the html page")
	}

	links := ProductLinks{BaseLink: link}
	logger.Println("Wil start processing node to get links to all products pages")
	err = links.ProcessNode(node, "Vezi toate produsele")
	if err != nil {
		logger.Fatalf("Could not process html node")
	}

	logger.Println("Got the following links:")
	for _, l := range links.links {
		logger.Println(l)
	}
}
