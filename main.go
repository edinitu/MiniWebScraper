package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/tebeka/selenium"
)

const (
	ERROR         = "ERROR: "
	CHROME        = "chrome"
	SELENIUM_PORT = 4444
)

var logger = log.Logger{}
var wd selenium.WebDriver

func InitSelenium() {
	var err error
	capabilities := selenium.Capabilities{"browserName": CHROME}
	wd, err = selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", SELENIUM_PORT))
	if err != nil {
		log.Fatalf("Failed to connect to the WebDriver: %v", err)
	}
}

func InitLogging() {
	logger.SetFlags(log.Ldate | log.Ltime)
	logger.SetOutput(os.Stdout)
}

func main() {
	InitLogging()
	InitSelenium()

	const link = "https://www.sephora.ro/"

	node, err := LinkToHtmlNode(link, false)

	if err != nil {
		logger.Fatalf(ERROR + "Could not get html node")
	}

	links := &ProductLinks{BaseLink: link}
	logger.Println("Wil start processing node to get links to all products pages")
	err = links.ProcessNode(node, "Vezi toate produsele")
	if err != nil {
		logger.Fatalf("Could not process html node")
	}

	logger.Printf("Got %d links: ", len(links.links))

	// TODO this should be made more abstract
	NewNode, err := LinkToHtmlNode(links.links[0], true)
	shop := &Shop{id: 1, name: "Sephora"}
	r := regexp.MustCompile("[a-zA-Z ]+\n[a-zA-Z ;+,]+\n[a-zA-Z ;+,]+\\n")
	err = shop.ProcessNode(NewNode, "")
	err = shop.ExtractProductsFromText([]*regexp.Regexp{r})
	// after one processing is done, clear global variable
	allText = ""
	if err != nil {
		logger.Fatalf("Couldn't process node for link %v", links.links[0])
	}

}
