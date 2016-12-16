package parser

import (
	"os"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	read()
}

func TestDocument_ReadDate(t *testing.T) {
	document, _ := read()

	date, _ := document.ReadDate()
	exp := time.Date(2010, 4, 1, 17, 0, 0, 0, date.Location())
	if !date.Equal(exp) {
		t.Errorf("Date is not equal. exp: '%v', got: '%v'", exp, date)
	}
}

func TestDocument_GetShop(t *testing.T) {
	document, _ := read()

	shop := document.GetShop()
	if shop == nil {
		t.Error("Shop cannot be nil")
	}

	if shop.Name == nil || *shop.Name != "shop name" {
		t.Errorf("Wrong shop name. exp: '%v', got: '%v'", "shop name", *shop.Name)
	}

	if shop.Company == nil || *shop.Company != "shop company" {
		t.Errorf("Wrong company. exp: '%v', got: '%v'", "shop company", *shop.Name)
	}

	if shop.URL == nil || *shop.URL != "shop url" {
		t.Errorf("Wrong url. exp: '%v', got: '%v'", "shop url", *shop.Name)
	}
}

func TestDocument_GetCategories(t *testing.T) {
	document, _ := read()
	shop := document.GetShop()
	if shop == nil {
		t.Error("Shop cannot be nil")
	}

	categories := shop.GetCategories()
	if categories.Length() != 19 {
		t.Errorf("GetCategories.Length(). exp: '%v', got: '%v'", 19, categories.Length())
	}

	category, found := categories.Get(1)
	if !found {
		t.Error("Category not found")
	}

	if category.Name != "Оргтехника" {
		t.Errorf("Wrong category name")
	}

}

func TestShop_ReadOffers(t *testing.T) {
	document, _ := read()
	shop := document.GetShop()
	if shop == nil {
		t.Error("Shop cannot be nil")
	}

	iter := shop.ReadOffers()
	list := make([]*Offer, 0)
	for {
		offer := iter()
		if offer == nil {
			break
		}

		list = append(list, offer)
	}

	if len(list) != 7 {
		t.Errorf("Wrong number of offers. exp: '%v', got: '%v'", 7, len(list))
	}

	table := map[string]Offer{
		"12341": {ID: "12341", Type:OffertTypeVendorModel, Available:true},
	}

	for _, item := range list {
		exp, found := table[item.ID]
		if !found {
			//TODO: enable me
			//t.Errorf("Offer not found, ID: %v", item.ID)
			continue
		}

		if *item != exp {
			t.Errorf("Offer is not equal. exp: '%+v', got: '%+v'", exp, item)
		}
	}
}

func read() (Document, error) {
	file, _ := os.Open("../fixtures/example.xml")

	return Read(file)
}
