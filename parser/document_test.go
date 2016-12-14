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

func read() (Document, error) {
	file, _ := os.Open("../fixtures/example.xml")

	return Read(file)
}
