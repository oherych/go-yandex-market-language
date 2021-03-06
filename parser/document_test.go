package parser

import (
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	offersTable = map[string]Offer{
		"12341": {
			ID:        "12341",
			Type:      OffertTypeVendorModel,
			Available: true,
			Vendor:    "НP",
			Model:     "Color LaserJet 3000",
			URL:       "http://magazin.ru/product_page.asp?pid=14344",
			Name:      "Принтер НP Color LaserJet 3000",
			Picture:   []string{"http://magazin.ru/img/device14344.jpg"},
			Price:     "15000",
			OldPrice:  "25000",
			Param: map[string]string{
				"загрузка": "100",
				"скорость": "3",
			},
			CurrencyID: "RUR",
			CategoryID: "101",
		},
	}

	categoriesTable = map[uint]Category{
		1: {},
	}
)


func TestShop_ReadOffers(t *testing.T) {
	file, _ := os.Open("../fixtures/example.xml")
	document, _ := Read(file)
	date, _ := document.ReadDate()
	exp := time.Date(2010, 4, 1, 17, 0, 0, 0, date.Location())
	if !date.Equal(exp) {
		t.Errorf("Date is not equal. exp: '%v', got: '%v'", exp, date)
	}

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

	categories := shop.GetCategories()
	if categories.Length() != 19 {
		t.Errorf("GetCategories.Length(). exp: '%v', got: '%v'", 19, categories.Length())
	}

	category, found := categories.Get(1)
	if !found {
		t.Error("Category not found")
	}

	if category.Name != "Оргтехника" {
		t.Error("Wrong category name")
	}
	
	categoriesList := categories.GetAll()
	category, found = categoriesList[1]
	if !found {
		t.Error("Category not found")
	}

	if category.Name != "Оргтехника" {
		t.Error("Wrong category name")
	}

	currencies := shop.GetCurrencies()
	if currencies.Length() != 5 {
		t.Errorf("GetCategories.Length(). exp: '%v', got: '%v'", 5, currencies.Length())
	}

	currency, found := currencies.Get("USD")
	if !found {
		t.Error("Currency not found")
	}

	if currency.rate != 23.98 {
		t.Error("Wrong currency rate")
	}

	iter := shop.ReadOffers()
	list := make([]Offer, 0)
	for {
		offer := iter()
		if offer == nil {
			break
		}

		list = append(list, *offer)
	}

	if len(list) != 7 {
		t.Errorf("Wrong number of offers. exp: '%v', got: '%v'", 7, len(list))
	}

	for _, item := range list {
		exp, found := offersTable[item.ID]
		if !found {
			//TODO: enable me
			//t.Errorf("Offer not found, ID: %v", item.ID)
			continue
		}

		if !reflect.DeepEqual(item, exp) {
			t.Errorf("Offer is not equal. exp: '%+v', got: '%+v'", exp, item)
		}
	}
}
