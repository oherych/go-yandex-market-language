package parser

import (
	xmlpath "gopkg.in/xmlpath.v2"
	"strconv"
)

const (
	CurrencyRateCBRF CurrencyRate = -1
	CurrencyRateNBU  CurrencyRate = -2
	CurrencyRateNBK  CurrencyRate = -3
	CurrencyRateСВ   CurrencyRate = -4
)

type CurrencyRate float32

type Currency struct {
	ID   string
	rate CurrencyRate
}

func NewCurrencies() (c Currencies) {
	c.list = make(map[string]Currency, 0)
	return
}

type Currencies struct {
	list map[string]Currency
}

func (c *Currencies) Add(currency Currency) {
	c.list[currency.ID] = currency
}

func (c Currencies) Get(ID string) (Currency, bool) {
	element, found := c.list[ID]
	return element, found
}

func (c Currencies) Length() int {
	return len(c.list)
}

func (c *Currencies) load(root *xmlpath.Node) {
	iter := xmlpath.MustCompile("currencies/currency").Iter(root)
	for iter.Next() {
		currency := Currency{}

		node := iter.Node()
		if val, ok := xmlpath.MustCompile("@id").String(node); ok {
			currency.ID = val
		} else {
			//TODO: ERROR
			continue
		}

		if val, ok := xmlpath.MustCompile("@rate").String(node); ok {
			rate, err := strconv.ParseFloat(val, 32)
			if err != nil {
				//TODO: ERROR
				continue
			}

			currency.rate = CurrencyRate(rate)
		}

		c.Add(currency)
	}
}
