package cache

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"sync"
)

type Cache struct {
	filePath string
	cache map[string]string
	lock sync.Mutex
}

func NewCache(filePath string) (*Cache, error) {
	var ret map[string]string
	cache := &Cache{
		filePath: filePath,
		cache:    nil,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		return cache, nil
	}

	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(body)
	decoder := gob.NewDecoder(buff)

	err = decoder.Decode(&ret)
	cache.cache = ret
	if err != nil {
		return cache, err
	}

	return cache, nil
}

func (cache *Cache) write() error {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(cache.cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cache.filePath, buff.Bytes(), 0666)
	if err != nil {
		return err
	}

	return nil
}

func (cache *Cache) NeedUpdate(ip string, ipType string) (bool, error) {
	if cache.cache == nil {
		cache.cache = map[string]string{
			ipType: ip,
		}

		cache.lock.Lock()
		err := cache.write()
		cache.lock.Unlock()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if oldIpAddr, ok := cache.cache[ipType]; ok {
		if oldIpAddr != ip {
			cache.cache[ipType] = ip

			cache.lock.Lock()
			err := cache.write()
			cache.lock.Unlock()
			if err != nil {
				return false, err
			}

			return true, nil
		}

		return false, nil
	} else {
		cache.cache[ipType] = ip

		cache.lock.Lock()
		err := cache.write()
		cache.lock.Unlock()
		if err != nil {
			return false, err
		}

		return true, nil
	}
}
