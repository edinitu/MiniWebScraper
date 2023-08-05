package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type HtmlExtractorInterface interface {
	GetHtmlPage(selenium bool) (string, error)
}

type HtmlExtractor struct {
	HtmlPage string
	link     string
}

type htmlRecord []byte

func (h *htmlRecord) Write(w []byte) (int, error) {
	for _, b := range w {
		*h = append(*h, b)
	}
	return len(w), nil
}

func (h *HtmlExtractor) GetHtmlPage(selenium bool) (string, error) {
	if !selenium {
		request, err := BuildHttpRequest(h.link)

		if err != nil {
			logger.Println(ERROR + "Could not get the html page")
			return "", err
		}

		str, err := SendRequestGetHtmlBody(request)
		if err != nil {
			logger.Println(ERROR + "Could not get the html page")
			return "", err
		}
		h.HtmlPage = str
		return str, nil
	}
	h.HtmlPage = GetPageWithSelenium(h.link)
	return h.HtmlPage, nil
}

func GetPageWithSelenium(link string) string {
	defer wd.Quit()

	// Navigate to the target website
	if err := wd.Get(link); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	// Wait for dynamic content to load (you may need to adjust the wait time based on the website)
	time.Sleep(5 * time.Second)

	// Get the dynamic HTML content
	// TODO make scrolls a configuration per shop, should be >= 1
	var scrolls int = 16
	var pageSource string = ""

	logger.Printf("Scrolling the page %v", link)
	for scrolls > 0 {
		_, err := wd.ExecuteScript("window.scrollTo(0,document.body.scrollHeight);", nil)
		if err != nil {
			log.Fatalf("Failed to scroll source page: %v", err)
		}
		// wait for the new elements to load
		time.Sleep(2 * time.Second)
		page, _ := wd.PageSource()
		pageSource += page
		scrolls--
	}

	return pageSource
}

// @ToBeTested
func BuildHttpRequest(link string) (*http.Request, error) {
	logger.Println("Sending GET request to", link)
	request, err := http.NewRequest("GET", link, nil)

	if err != nil {
		logger.Println(ERROR + "Could not send request to site")
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("User-Agent", "Mozilla/5.0")

	return request, nil
}

func SendRequestGetHtmlBody(request *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Println(ERROR+"Could not get response from server ", err)
		return "", err
	}

	logger.Println("Got response status:", resp.Status)

	hr := htmlRecord{}
	_, err = io.Copy(&hr, resp.Body)
	if err != nil {
		logger.Println(ERROR+"Could not parse response body into string ", err)
		return "", err
	}
	return string(hr), nil
}

/*
Get a node to the root HTML element from the HTML page given as a string.
@ ToBeTested
*/
func ParseHtml(text string) (*html.Node, error) {
	r := strings.NewReader(text)
	node, err := html.Parse(r)
	if err != nil {
		logger.Println(ERROR+" Could not parse html page ", err)
		return nil, err
	}

	return node, nil
}

/*
Given a link to a webpage, return the head HTML node for that page.

	Arguments: - the webpage link as a string
			   - flag for selenium use in case we want to load JS elements that are
				 not loaded automatically from the usual http.Request
*/
func LinkToHtmlNode(link string, selenium bool) (*html.Node, error) {
	extractor := HtmlExtractor{link: link}
	hp, err := extractor.GetHtmlPage(selenium)
	if err != nil {
		return nil, err
	}

	node, err := ParseHtml(hp)
	if err != nil {
		logger.Println(ERROR, "Could not parse the html page")
		return nil, err
	}
	return node, nil
}
