// package main
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	// "fmt"	strings"
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

	categories []string

	// TODO regexs should be a field here
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
	logger.Println("Initialize products from text")
	firstMatches := r[0].FindAllString(allText, -1)
	filteredValues := firstFilter(firstMatches)
	filteredValues = secondFilter(filteredValues)

	//	products := initProducts(f
	products := initProducts(filteredValues, sh.id)
	setRemainingInfo(products, []*regexp.Regexp{r[1], r[2]})
	//fmt.Println(filteredValues)
	fmt.Println(filteredValues)
	return nil
}

// @ToBeTested
func initProducts(s []string, shopId uint64) map[string]Product {
	var c uint64 = 0
	products := make(map[string]Product)
	// TODO This should be based on shop ID
	for _, seq := range s {
		lines := strings.Split(seq, "\n")
		p := Product{
			id:       c,
			ShopId:   shopId,
			brand:    strings.TrimSpace(lines[0]),
			name:     strings.TrimSpace(lines[1]),
			category: strings.TrimSpace(lines[2]),
		}
		_, ok := products[p.name]
		if !ok {
			products[p.name] = p
		}
	}
	return products
}

func setRemainingInfo(products map[string]Product, rl []*regexp.Regexp) {
	for key, p := range products {
		r := regexp.MustCompile(p.name)
		s := r.FindIndex([]byte(allText))
		p.price = rl[0].FindString(allText[s[1] : s[1]+80])
		PricePerQuantity := rl[1].FindString(allText[s[1] : s[1]+100])
		PricePerQuantity = strings.Replace(PricePerQuantity, "\n", "", -1)
		// TODO Change the following assignment when getQuantity() is done
		p.quantity = strings.TrimSpace(PricePerQuantity)
		fmt.Println(p.GetStringRepresentation())
		products[key] = p
		os.Exit(1)
	}
}

func getQuantity(price string, PricePerQuantity string) string {
	//TODO
	return ""
}

// @ToBeTested
func firstFilter(s []string) []string {
	res := []string{}
	for _, seq := range s {
		if len(seq) > 5 && isTruePositive(seq) {
			res = append(res, seq)
		}
	}
	return res
}

// @ToBeTested
func secondFilter(s []string) []string {
	res := []string{}
	for _, seq := range s {
		firstLine := seq[:strings.IndexByte(seq, byte('\n'))]
		if isJustUpperCase(firstLine) {
			res = append(res, seq)
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

// @ToBeTested
func isJustUpperCase(s string) bool {
	for _, r := range s {
		if !(unicode.IsSpace(r) || unicode.IsUpper(r)) {
			return false
		}
	}
	return true
}
