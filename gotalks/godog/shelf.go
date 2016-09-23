package main

// Shelf stores a list of products which are available for purchase type Shelf
type Shelf struct {
	products map[string]int
}

func (s *Shelf) Add(title string, price int) {
	if nil == s.products {
		panic("shelf products are not initialized")
	}

	s.products[title] = price
}
