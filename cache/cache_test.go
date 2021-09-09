package cache

import (
	"testing"
	"time"
)

func TestStoreTweetsFromStream(t *testing.T) {
	testCache := NewCache()
	idList := [11]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}

	for i := 0; i < len(idList); i++ {
		go testCache.StoreTweetFromStream(idList[i])
		time.Sleep(time.Millisecond)
	}
	stateOfCache, found := testCache.cache.Get(COLLECTED_TWEETS); if !found {
		t.Errorf(`unable to find list in cache`)
	}

	cl := stateOfCache.([]string)

	if len(cl) != 1 {
		t.Errorf(`cache items still persist in cache after flush, number of items: %d`, len(cl))
	}

}