package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type NodeProcessor interface {
	ProcessNode(n *html.Node, s string) error
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

func (sh *Shop) ProcessNode(n *html.Node, str string) error {
	if n.Type == html.ElementNode && n.DataAtom == atom.Span {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			fmt.Println(n.FirstChild.Data)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sh.ProcessNode(c, str)
	}
	return nil
}
