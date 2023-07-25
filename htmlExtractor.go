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
	pageSource, err := wd.PageSource()
	if err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}
	return pageSource
}

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
