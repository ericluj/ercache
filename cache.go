package ercache

import (
	"encoding/json"
	"io"
	"sync"
)

type cache struct {
	data map[string]string
	sync.RWMutex
}

func newCache() *cache {
	return &cache{
		data: make(map[string]string),
	}
}

func (c *cache) Get(key string) string {
	c.RLock()
	defer c.RUnlock()

	if val, ok := c.data[key]; ok {
		return val
	}
	return ""
}

func (c *cache) Set(key, value string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
	return nil
}

func (c *cache) Marshal() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	dataBytes, err := json.Marshal(c.data)
	return dataBytes, err
}

func (c *cache) Unmarshal(serialized io.ReadCloser) error {
	var newData map[string]string
	if err := json.NewDecoder(serialized).Decode(&newData); err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()
	c.data = newData
	return nil
}
