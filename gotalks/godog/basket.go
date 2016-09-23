package main

type Basket struct {
	products map[string]int
}

func (b *Basket) Add(title string, price int) {
	if nil == b.products {
		panic("basket products are not initialized")
	}

	b.products[title] = price
}

func (b *Basket) Total() int {
	var total int
	for _, value := range b.products {
		total += value
	}

	total = int(float64(total) * 1.2)
	if total <= 10 {
		return total + 3
	}

	return total + 2
}
