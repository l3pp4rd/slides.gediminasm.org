package main

import (
	"fmt"

	"github.com/DATA-DOG/godog"
)

type shoppingContext struct {
	shelf  *Shelf
	basket *Basket
}

func (sc *shoppingContext) thereIsAProductWhichCosts(title string, price int) error {
	sc.shelf.Add(title, price)
	return nil
}

func (sc *shoppingContext) iAddProductToTheBasket(title string) error {
	if price, exists := sc.shelf.products[title]; !exists {
		return fmt.Errorf("product %s could not be found on shelf", title)
	} else {
		sc.basket.products[title] = price
	}
	return nil
}

func (sc *shoppingContext) iShouldHaveNumProductsInTheBasket(num int) error {
	if num != len(sc.basket.products) {
		return fmt.Errorf("basket has %d items", len(sc.basket.products))
	}
	return nil
}

func (sc *shoppingContext) theOverallBasketPriceShouldBe(price int) error {
	if sc.basket.Total() != price {
		return fmt.Errorf("actual total price is: %d", sc.basket.Total())
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	var context shoppingContext

	s.BeforeScenario(func(_ interface{}) {
		context.shelf = &Shelf{products: make(map[string]int)}
		context.basket = &Basket{products: make(map[string]int)}
	})

	s.Step(`^there is a "([^"]*)", which costs £(\d+)$`, context.thereIsAProductWhichCosts)
	s.Step(`^I add the "([^"]*)" to the basket$`, context.iAddProductToTheBasket)
	s.Step(`^I should have (\d+) products? in the basket$`, context.iShouldHaveNumProductsInTheBasket)
	s.Step(`^the overall basket price should be £(\d+)$`, context.theOverallBasketPriceShouldBe)
}
