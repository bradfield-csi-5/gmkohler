package pkg

import "strings"

type responses map[string][]byte
type HttpCache struct {
	CacheConfig
	cached responses
}

func NewHttpCache(config CacheConfig) *HttpCache {
	return &HttpCache{
		CacheConfig: config,
		cached:      make(responses),
	}
}

func (c *HttpCache) ShouldCache(path string) bool {
	return strings.HasPrefix(path, c.Prefix)
}

func (c *HttpCache) Get(path string) []byte {
	data, ok := c.cached[path]
	if !ok {
		return nil
	}
	return data
}

// CacheConfig - Write logic that allows an administrator to specify a path,
// and have your proxy automatically cache any request to a URL
// prefixed by that path
type CacheConfig struct {
	Prefix string
}
