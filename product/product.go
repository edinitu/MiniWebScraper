package product

import "fmt"

type Product struct {
	Id            uint64
	ShopId        uint64
	Brand         string
	Name          string
	Category      string
	Price         string
	IsPromotion   bool
	OriginalPrice string
	Quantity      string
}

func (p *Product) SetDefaults() {
	if !p.IsPromotion {
		p.OriginalPrice = p.Price
	}
}

func (p Product) CheckValues() {
	//TODO
}

func (p Product) GetStringRepresentation() string {
	return "Product description: \n" +
		"	ID: " + fmt.Sprint(p.Id) + "\n" +
		"	ShopId: " + fmt.Sprint(p.ShopId) + "\n" +
		"	Brand: " + p.Brand + "\n" +
		"	Name: " + p.Name + "\n" +
		"	Category: " + p.Category + "\n" +
		"	Price: " + p.Price + "\n" +
		"	Promotion: " + fmt.Sprint(p.IsPromotion) + "\n" +
		"	Original price: " + p.OriginalPrice + "\n" +
		"	Quantity: " + p.Quantity
}

func (p Product) Serialize() []byte {
	//TODO, JSON??
	return nil
}
