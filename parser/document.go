package parser

import (
	"encoding/xml"
	"errors"
	"gopkg.in/xmlpath.v2"
	"io"
	"time"
)

type Document struct {
	root *xmlpath.Node
}

type Shop struct {
	root *xmlpath.Node

	Name    *string
	Company *string
	URL     *string
}

func Read(r io.Reader) (Document, error) {
	xmlDecoder := xml.NewDecoder(r)
	root, err := xmlpath.ParseDecoder(xmlDecoder)
	if err != nil {
		return Document{}, err
	}

	return Document{root: root}, nil
}

func (d Document) GetShop() *Shop {
	iter := xmlpath.MustCompile("/yml_catalog/shop").Iter(d.root)
	if !iter.Next() {
		return nil
	}

	return newShop(iter.Node())
}

func (d Document) ReadDate() (time.Time, error) {
	if val, ok := xmlpath.MustCompile("/yml_catalog/@date").String(d.root); ok {
		return time.Parse("2006-01-02 15:04", val)
	}

	return time.Now(), errors.New("Cannot find date")
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
