package dig

import (
	"bufio"
	"hash"
	"io"
	"os"
	"sync"

	"golang.org/x/sync/singleflight"
)

type fileHashCacheItem struct {
	info os.FileInfo
	hash []byte
}

type fileHashCache struct {
	mu    sync.RWMutex
	items map[string]*fileHashCacheItem

	sf singleflight.Group
}
var fileHashCacheInstance = &fileHashCache{}

func (c *fileHashCache) Get(filepath string, hsh hash.Hash) ([]byte, error) {
	res, err, _ := c.sf.Do(filepath, func() (any, error) {
		var item *fileHashCacheItem
		c.mu.RLock()
		if c.items != nil {
			item = c.items[filepath]
		}
		c.mu.RUnlock()
		var info os.FileInfo
		if item != nil {
			var err error
			info, err = os.Stat(filepath)
			if err != nil {
				return nil, err
			}
			if info.ModTime().Equal(item.info.ModTime()) &&
				info.Size() == item.info.Size() {
				return item.hash, nil
			}
		}
		// If the file is not in the cache or has changed, calculate the hash.
		hash, err := hashFile(filepath, hsh)
		if err != nil {
			return nil, err
		}
		c.mu.Lock()
		if c.items == nil {
			c.items = make(map[string]*fileHashCacheItem)
		}
		c.items[filepath] = &fileHashCacheItem{
			info: info,
			hash: hash,
		}
		c.mu.Unlock()
		return hash, nil
	})
	if err != nil {
		return nil, err
	}
	return res.([]byte), nil
}

func hashFile(filepath string, hsh hash.Hash) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReaderSize(f, 64*1024*1024)
	_, err = io.Copy(hsh, r)
	if err != nil {
		return nil, err
	}
	return hsh.Sum(nil), nil
}
