package cache
import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

const COLLECTED_TWEETS = "collectedtweets"

type Cache struct {
	cache *cache.Cache
	mutex  *sync.Mutex
}

func NewCache() *Cache {
	cache := cache.New(60*time.Minute, 70*time.Minute)
	var emptyIdList []string
	cache.Set(COLLECTED_TWEETS, emptyIdList, 60*time.Minute)
	mutex := &sync.Mutex{}
	return &Cache{
		cache,
		mutex,
	}
}

func (c *Cache) GetCache(key string) interface{} {
	ct, found := c.cache.Get(COLLECTED_TWEETS); if !found {
		return nil
	}
	return ct.([]string)
}

func (c *Cache) StoreTweetFromStream(id string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ct, found := c.cache.Get(COLLECTED_TWEETS); if !found {
		fmt.Println(fmt.Sprintf(`rebuilding cache from id: %s`, id))
		el := []string{id}
		c.cache.Set(COLLECTED_TWEETS, el, 60*time.Minute)
		return nil
	}
	fmt.Println(fmt.Sprintf(`adding new element to cache %s`, id))
	newList := append(ct.([]string), id)
	c.cache.Set(COLLECTED_TWEETS, newList, 60*time.Minute)
	cs, err := c.CheckCacheSize(); if err != nil {
		return err
	}

	if cs {
		c.ProcessTweetsFromCache()
	}
	return nil
}

func (c *Cache) CheckCacheSize() (bool, error) {
	ct, found := c.cache.Get(COLLECTED_TWEETS); if !found {
		return false, fmt.Errorf("unable to find collected tweets")
	}
	// check cache size
	cacheSize := len(ct.([]string))
	if cacheSize >= 10 {
		return true, nil
	}
	return false, nil
}

func (c *Cache) ProcessTweetsFromCache() {
	cacheState := c.GetCache(COLLECTED_TWEETS)
	fmt.Println(fmt.Sprintf(`cache state before flush %v` , cacheState))
	// get info for all tweets currently in cache

	c.cache.Flush()
}