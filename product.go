package main

import "fmt"

type Product struct {
	id            uint64
	ShopId        uint64
	brand         string
	name          string
	category      string
	price         string
	IsPromotion   bool
	OriginalPrice string
	quantity      string
}

func (p *Product) SetDefaults() {
	if !p.IsPromotion {
		p.OriginalPrice = p.price
	}
}

func (p Product) CheckValues() {
	//TODO
}

func (p Product) GetStringRepresentation() string {
	return "Product description: \n" +
		"	ID: " + fmt.Sprint(p.id) + "\n" +
		"	ShopId: " + fmt.Sprint(p.ShopId) + "\n" +
		"	Brand: " + p.brand + "\n" +
		"	Name: " + p.name + "\n" +
		"	Category: " + p.category + "\n" +
		"	Price: " + p.price + "\n" +
		"	Promotion: " + fmt.Sprint(p.IsPromotion) + "\n" +
		"	Original price: " + p.OriginalPrice + "\n" +
		"	Quantity: " + p.quantity
}

func (p Product) Serialize() []byte {
	//TODO, JSON??
	return nil
}
