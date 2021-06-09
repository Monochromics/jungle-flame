package main

import (
	"github.com/dgraph-io/ristretto"
)

var c *ristretto.Cache

func cacheInit() (cache *ristretto.Cache) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return
}
