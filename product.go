package main

type Product struct {
	id            uint64
	ShopId        uint64
	name          string
	category      string
	price         float64
	IsPromotion   bool
	OriginalPrice float64
	quantity      int64
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
	//TODO
	return ""
}

func (p Product) Serialize() []byte {
	//TODO, JSON??
	return nil
}
