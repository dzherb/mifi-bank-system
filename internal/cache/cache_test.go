package cache_test

import (
	"testing"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/cache"
)

func TestCache(t *testing.T) {
	cache.Init()
	t.Cleanup(cache.Close)

	const sleep = time.Millisecond * 100

	cache.Active().Set("test", "test_val", sleep)

	item := cache.Active().Get("test")
	if v := item.Value().(string); v != "test_val" {
		t.Errorf("got %s, want test_val", v)
	}

	time.Sleep(sleep)

	item = cache.Active().Get("test")
	if item != nil {
		t.Errorf("got %s, want nil", item.Value().(string))
	}
}
