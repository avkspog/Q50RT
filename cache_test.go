package main

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := newCache()

	v, found := cache.Get(1234567890)
	if found || v != nil {
		t.Error("value should not be found!")
	}

	cache.Set(1234567890, "test_value_1")
	cache.Set(123456789, "test_value_2")

	v, found = cache.Get(1234567890)
	if !found || v == nil {
		t.Error("test_value_1 doesn't found!")
	}

	v, found = cache.Get(123456789)
	if !found || v == nil {
		t.Error("test_value_2 doesn't found!")
	}
}

func TestCacheTimeout(t *testing.T) {
	cache := newCache()

	cache.SetExp(123456789, "exp_1_sec", 100*time.Millisecond)

	<-time.After(50 * time.Millisecond)
	v, found := cache.Get(123456789)
	if !found || v == nil {
		t.Error("exp_1_sec is expired", v)
	}

	<-time.After(50 * time.Millisecond)
	v, found = cache.Get(123456789)
	if found || v != nil {
		t.Error("exp_1_sec doesn't expired", v)
	}

}
