package cache

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func setCache(count int, b *testing.B) {
	ctx := context.Background()

	lruCache := NewCache(count, time.Second*5)

	b.ResetTimer()
	for i := 0; i < count; i++ {
		lruCache.Set(ctx, strconv.Itoa(i), i)
	}
}

func setAndGetCache(count int, b *testing.B) {
	ctx := context.Background()

	lruCache := NewCache(count, time.Second*5)

	for i := 0; i < count; i++ {
		lruCache.Set(ctx, strconv.Itoa(i), i)
	}

	b.ResetTimer()
	for i := 0; i < count; i++ {
		lruCache.Get(ctx, strconv.Itoa(i))
	}
}

func BenchmarkCacheLookup1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setAndGetCache(1000, b)
	}
}

func BenchmarkCacheLookup10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setAndGetCache(10000, b)
	}
}

func BenchmarkCacheLookup20000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setAndGetCache(20000, b)
	}
}

func BenchmarkCacheMemory1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setCache(1000, b)
	}
}

func BenchmarkCacheMemory10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setCache(10000, b)
	}
}

func BenchmarkCacheMemory20000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setCache(20000, b)
	}
}
