// package main
package main

import (
	"go_prj/db"
	"go_prj/product"
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
	ProcessNode(n *html.Node, s string) error
}

type ShopService interface {
	ExtractProductsFromText() error
	SetCategories(p *ProductLinks)
}

type ProductLinks struct {
	BaseLink string
	links    []string
}

type Shop struct {
	id            uint64
	name          string
	patterns      []*regexp.Regexp
	CategoryLinks map[string]string
	db            *db.DbConn
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

/*
* Using the HTML text extracted, populate a map with all the products using different filters.
 */
func (sh *Shop) ExtractProductsFromText(category string) error {
	logger.Println("Initialize products from text")
	firstMatches := sh.patterns[0].FindAllString(allText, -1)
	filteredValues := firstFilter(firstMatches)
	filteredValues = secondFilter(filteredValues)

	products := initProducts(filteredValues, sh.id)
	setRemainingInfo(
		products,
		[]*regexp.Regexp{sh.patterns[1], sh.patterns[2], sh.patterns[3]},
		category)
	err := persistProducts(products, sh.db)
	if err != nil {
		logger.Printf(ERROR+"Could not persist products to DB, %v", err)
		return err
	}
	logger.Printf("Persisted %d products to DB", len(products))
	return nil
}

func (s *Shop) SetCategories(p *ProductLinks) {
	s.CategoryLinks = make(map[string]string)
	for _, link := range p.links {
		s.CategoryLinks[GetCategoryFromLink(link, s.patterns[4])] = link
	}
}

func GetCategoryFromLink(link string, r *regexp.Regexp) string {
	s := r.Find([]byte(link))
	re := strings.Trim(string(s), "-")
	return strings.Trim(re, "/")
}

// @ToBeTested
func initProducts(s []string, shopId uint64) map[string]product.Product {
	var c uint64 = 0
	products := make(map[string]product.Product)
	// TODO This should be based on shop ID
	for _, seq := range s {
		lines := strings.Split(seq, "\n")
		p := product.Product{
			Id:       c,
			ShopId:   shopId,
			Brand:    strings.TrimSpace(lines[0]),
			Name:     strings.TrimSpace(lines[1]),
			Category: strings.TrimSpace(lines[2]),
		}
		_, ok := products[p.Name]
		if !ok {
			products[p.Name] = p
		}
	}
	return products
}

func setRemainingInfo(products map[string]product.Product, rl []*regexp.Regexp, category string) {
	notFoundProducts := []string{}
	var count int = 0
	for key, p := range products {
		r := regexp.MustCompile(p.Name)
		s := r.FindIndex([]byte(allText))
		if s == nil {
			notFoundProducts = append(notFoundProducts, p.Name)
			continue
		}
		p.Price = rl[0].FindString(allText[s[1] : s[1]+80])
		if p.Price == "" {
			count++
		}
		PricePerQuantity := rl[1].FindString(allText[s[1] : s[1]+100])
		PricePerQuantity = strings.Replace(PricePerQuantity, "\n", "", -1)
		if s[0]-20 > 0 {
			if rl[2].FindString(allText[s[0]-60:s[0]]) != "" {
				p.IsPromotion = true
			}
		}
		// TODO Change the following assignment when getQuantity() is done
		p.Quantity = strings.TrimSpace(PricePerQuantity)
		p.Category = category
		p.SetDefaults()
		products[key] = p
	}

	if len(notFoundProducts) > 0 {
		logger.Println("Could not set remaining info for these products:", notFoundProducts)
	}
	if count > 0 {
		logger.Println("Could not set price for number of products:", count)
	}
}

func persistProducts(products map[string]product.Product, dbConn *db.DbConn) error {
	//TODO get PostgreSQL connection and bulk insert the received products
	err := dbConn.Db.Ping()
	if err != nil {
		return err
	}
	db.InsertProducts(products, dbConn)
	return nil
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
