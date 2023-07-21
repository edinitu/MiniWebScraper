package main

import (
	"io"
	"net/http"
)

type HtmlExtractorInterface interface {
	GetHtmlPage() string
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

func (h *HtmlExtractor) GetHtmlPage() (string, error) {
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
