package parser

import (
	"encoding/xml"
	"errors"
	"gopkg.in/xmlpath.v2"
	"io"
	"strconv"
	"time"
)

type offerInterator func() *Offer
type OfferType string

type Document struct {
	root *xmlpath.Node
}

type Shop struct {
	root *xmlpath.Node

	Name    *string
	Company *string
	URL     *string

	categories *Categories
}

type Category struct {
	ID       uint
	ParentID uint
	Name     string
}

type Categories struct {
	list map[uint]Category
}

type Offer struct {
	ID        string
	Type      OfferType
	Available bool
	Vendor    string
	Model     string
	URL string
}

const (
	OffertTypeDefault     OfferType = ""
	OffertTypeVendorModel OfferType = "vendor.model"
	OffertTypeMedicine    OfferType = "medicine"
	OffertTypeBooks       OfferType = "books"
	OffertTypeAudioBooks  OfferType = "audiobooks"
	OffertTypeArtistTitle OfferType = "artist.title"
	OffertTypeEventTicket OfferType = "event-ticket"
	OffertTypeTour        OfferType = "tour"
)

func Read(r io.Reader) (Document, error) {
	xmlDecoder := xml.NewDecoder(r)
	root, err := xmlpath.ParseDecoder(xmlDecoder)
	if err != nil {
		return Document{}, err
	}

	return Document{root: root}, nil
}

func NewCategories() (c Categories) {
	c.list = make(map[uint]Category, 0)
	return
}

func (o OfferType) String() string {
	return string(o)
}

func (d Document) GetShop() *Shop {
	iter := xmlpath.MustCompile("/yml_catalog/shop").Iter(d.root)
	if !iter.Next() {
		return nil
	}

	return newShop(iter.Node())
}

func (s *Shop) GetCategories() *Categories {
	if s.categories == nil {
		c := createCategories(s.root)
		s.categories = &c
	}

	return s.categories
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

func (o *Offer) LoadFromNode(node *xmlpath.Node) {
	if val, ok := xmlpath.MustCompile("@id").String(node); ok {
		o.ID = val
	}

	if val, ok := xmlpath.MustCompile("@type").String(node); ok {
		o.Type = OfferType(val)
	}

	if val, ok := xmlpath.MustCompile("@available").String(node); ok {
		o.Available = (val == "true")
	} else {
		o.Available = true
	}

	if val, ok := xmlpath.MustCompile("vendor").String(node); ok {
		o.Vendor = val
	}

	if val, ok := xmlpath.MustCompile("model").String(node); ok {
		o.Model = val
	}

	if val, ok := xmlpath.MustCompile("url").String(node); ok {
		o.URL = val
	}
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

func (c *Categories) Add(category Category) {
	c.list[category.ID] = category
}

func (c Categories) Length() int {
	return len(c.list)
}

func (c Categories) Get(ID uint) (Category, bool) {
	element, found := c.list[ID]
	return element, found
}

func newShop(root *xmlpath.Node) *Shop {
	shop := Shop{root: root}
	shop.load()
	return &shop
}

func createCategories(root *xmlpath.Node) Categories {
	c := NewCategories()

	iter := xmlpath.MustCompile("categories/category").Iter(root)
	for iter.Next() {
		category := Category{}

		node := iter.Node()
		if val, ok := xmlpath.MustCompile("@id").String(node); ok {
			id, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				//TODO: ERROR
				continue
			}

			category.ID = uint(id)
		}

		if val, ok := xmlpath.MustCompile("@parentId").String(node); ok {
			id, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				//TODO: ERROR
				continue
			}

			category.ParentID = uint(id)
		}

		category.Name = node.String()

		c.Add(category)
	}

	return c
}
