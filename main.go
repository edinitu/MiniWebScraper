package main

import (
	"database/sql"
	"fmt"
	"go_prj/db"
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
	logger.Println("Init Selenium Connection")
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

func InitDbConnection() *sql.DB {
	logger.Println("Init DB Connection")
	pgconfig := db.PgConfig{
		User:     "postgres",
		Password: "admin",
		Dbname:   "nssx",
		Sslmode:  "disable",
	}
	dbConn, err := db.InitDb(pgconfig)
	if err != nil {
		logger.Fatalf("Could not connect to DB, abort, %v", err)
	}
	return dbConn
}

func main() {
	InitLogging()
	InitSelenium()
	dbConn := InitDbConnection()

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

	logger.Printf("Got %d links", len(links.links))

	// TODO this should be made more abstract in a different function/service
	// TODO This should be implemented for all links
	r1 := regexp.MustCompile("[a-zA-Z ]+\n[a-zA-Z ;+,]+\n[a-zA-Z ;+,]+\\n")
	r2 := regexp.MustCompile(`[0-9,. ]+Lei`)
	r3 := regexp.MustCompile(`[0-9,. ]+Lei\s+\/\s+[0-9]+[a-z]`)
	r4 := regexp.MustCompile("Promo")
	shop := &Shop{
		id:       1,
		name:     "Sephora",
		patterns: []*regexp.Regexp{r1, r2, r3, r4},
	}
	for _, link := range links.links {
		defer wd.Quit()
		NewNode, err := LinkToHtmlNode(link, true)
		if err != nil {
			logger.Fatalf("Could not get node for link %v", links.links[0])
		}
		err = shop.ProcessNode(NewNode, "")
		if err != nil {
			logger.Fatalf("Could not process node for shop %v", shop.name)
		}
		err = shop.ExtractProductsFromText(dbConn)
		// after one processing is done, clear global variable
		allText = ""
		if err != nil {
			logger.Fatalf("Couldn't process node for link %v", links.links[0])
		}
	}
}
