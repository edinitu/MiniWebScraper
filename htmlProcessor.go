package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// This holds the all the text extracted from a HTML page.
// Its whole purpose is to be a starting point for various
// filters in order to extract relevant information.
var allText string

type NodeProcessor interface {
	ProcessNode(n *html.Node, exps []*regexp.Regexp) error
}

type ProductLinks struct {
	BaseLink string
	links    []string
}

func (p *ProductLinks) ProcessNode(n *html.Node, s string) error {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			if strings.Contains(n.FirstChild.Data, s) {
				new_link := p.BaseLink + n.Attr[1].Val
				p.links = append(p.links, new_link)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.ProcessNode(c, s)
	}
	return nil
}

type Shop struct {
	id   uint64
	name string
}

func (sh *Shop) ProcessNode(n *html.Node, s string) error {
	if n.Type == html.ElementNode && n.DataAtom == atom.Span {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			allText += n.FirstChild.Data + "\n"
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sh.ProcessNode(c, s)
	}
	return nil
}

// TODO this should have a loop over multiple regex (to get info like promotions, price etc)
func (sh *Shop) ExtractProductsFromText(r []*regexp.Regexp) error {
	firstMatches := r[0].FindAllString(allText, -1)
	filteredValues := firstFilter(firstMatches)
	fmt.Println(filteredValues)
	return nil
}

// @ToBeTested
func firstFilter(s []string) []string {
	res := []string{}
	for _, w := range s {
		if len(w) > 5 && isTruePositive(w) {
			res = append(res, w)
		}
	}
	return res
}

// @ToBeTested
func isTruePositive(s string) bool {
	count := 0
	for _, c := range s {
		if unicode.IsLetter(c) {
			count++
		}
	}
	return count > 10
}
