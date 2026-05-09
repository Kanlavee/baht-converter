package main

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"

	"github.com/Kanlavee/baht-converter/internal/converter"
)

func main() {
	c := converter.NewConverter()

	inputs := []decimal.Decimal{
		decimal.NewFromFloat(1234),
		decimal.NewFromFloat(33333.75),
		decimal.NewFromFloat(0),
		decimal.NewFromFloat(11.01),
		decimal.NewFromFloat(21.50),
		decimal.NewFromFloat(1000001),
		decimal.NewFromFloat(-51),
		decimal.NewFromFloat(-0.5),
		decimal.NewFromFloat(0.0099999999),
		decimal.NewFromFloat(1.999),
		decimal.NewFromFloat(9999999.99),
		
	}

	for _, input := range inputs {
		result, err := c.Convert(input)
		if err != nil {
			log.Printf("Convert(%v) error: %v", input, err)
			continue
		}
		fmt.Printf("Input: %-15v -> Output: %s\n", input, result)
	}
}

