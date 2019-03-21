package namespace

import (
	"github.com/coreos/clair/database"
	"github.com/hashicorp/golang-lru"
)

type namespaceCache struct {
	c *lru.Cache
}

func newCache(size int) *namespaceCache {
	c, err := lru.New(size)
	if err != nil {
		panic(err)
	}

	return &namespaceCache{c}
}

func (c *namespaceCache) Get(ns database.Namespace) (int, bool) {
	if id, ok := c.c.Get(ns); ok {
		return id.(int), true
	}

	return 0, false
}

func (c *namespaceCache) Add(ns database.Namespace, id int) bool {
	return c.c.Add(ns, id) && c.c.Add(id, ns)
}

func (c *namespaceCache) GetValue(id int) (database.Namespace, bool) {
	if v, ok := c.c.Get(id); ok {
		return v.(database.Namespace), true
	}

	return database.Namespace{}, false
}
