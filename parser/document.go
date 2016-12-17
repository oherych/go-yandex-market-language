package parser

import (
	"encoding/xml"
	"errors"
	"io"
	"time"

	"gopkg.in/xmlpath.v2"
)

type Document struct {
	root *xmlpath.Node
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
