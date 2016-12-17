package parser

import xmlpath "gopkg.in/xmlpath.v2"

type Shop struct {
	root *xmlpath.Node

	Name    *string
	Company *string
	URL     *string

	categories *Categories
	currencies *Currencies
}

func (s *Shop) GetCategories() *Categories {
	if s.categories == nil {
		c := NewCategories()
		s.categories = &c
		s.categories.load(s.root)
	}

	return s.categories
}

func (s *Shop) GetCurrencies() *Currencies {
	if s.currencies == nil {
		c := NewCurrencies()
		s.currencies = &c
		s.currencies.load(s.root)
	}

	return s.currencies
}

func (s Shop) ReadOffers() offerInterator {
	iter := xmlpath.MustCompile("offers/offer").Iter(s.root)

	return func() *Offer {
		if !iter.Next() {
			return nil
		}

		offer := Offer{}
		offer.LoadFromNode(iter.Node())

		return &offer
	}
}

func (s *Shop) load() {
	if val, ok := xmlpath.MustCompile("name").String(s.root); ok {
		s.Name = &val
	}

	if val, ok := xmlpath.MustCompile("company").String(s.root); ok {
		s.Company = &val
	}

	if val, ok := xmlpath.MustCompile("url").String(s.root); ok {
		s.URL = &val
	}
}

func newShop(root *xmlpath.Node) *Shop {
	shop := Shop{root: root}
	shop.load()
	return &shop
}
