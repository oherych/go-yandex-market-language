package parser

import (
	"strconv"

	xmlpath "gopkg.in/xmlpath.v2"
)

type Category struct {
	ID       uint
	ParentID uint
	Name     string
}

type Categories struct {
	list map[uint]Category
}

func NewCategories() (c Categories) {
	c.list = make(map[uint]Category, 0)
	return
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

func (c *Categories) load(root *xmlpath.Node) {

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
}
