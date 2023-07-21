package main

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type NodeProcessor interface {
	ProcessNode(n *html.Node, s string) (interface{}, error)
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
