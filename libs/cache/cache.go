package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

const (
	// константы для хэш-функции
	// https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function#FNV-1a_hash
	fnvOffSetBasis uint64 = 14695981039346656037
	fnvPrime       uint64 = 1099511628211
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) bool
	Get(ctx context.Context, key string) interface{}
}

// элемент кэша
type cacheItem struct {
	key   string
	value interface{}

	ttl *time.Time
}

type bucket struct {
	// ключ - значение
	data map[string]*list.Element
	// для thread safe
	mu *sync.Mutex
}

// кэш
type cache struct {
	// ёмкость кэша
	capacity int
	// слайс бакетов
	buckets []*bucket
	//очередь для LRU
	queue *list.List
	// тикер для очистки кэша по ttl
	ticker *time.Ticker
	// ttl по умолчанию
	ttl time.Duration
}

// конструктор
func NewCache(capacity int, ttl time.Duration) Cache {
	buckets := make([]*bucket, 16)
	for index := range buckets {
		buckets[index] = &bucket{
			data: make(map[string]*list.Element),
			mu:   &sync.Mutex{},
		}
	}

	cache := &cache{
		capacity: capacity,
		buckets:  buckets,
		queue:    list.New(),
		ticker:   time.NewTicker(time.Second * 10),
		ttl:      ttl,
	}

	// запускаем горутину для очистки кэша по ttl
	go cache.cleanCache()

	return cache
}

// хэш-функция для определения номера бакета
func hashValue(key string, limit int) int {
	hash := fnvOffSetBasis
	for _, b := range []byte(key) {
		hash ^= uint64(b)
		hash *= fnvPrime
	}
	return int(hash % uint64(limit))
}

// записываем значение в кэш
func (c *cache) Set(ctx context.Context, key string, value interface{}) bool {
	// определяем индекс бакета
	hash := hashValue(key, len(c.buckets))

	// блокировка
	c.buckets[hash].mu.Lock()
	defer c.buckets[hash].mu.Unlock()

	// проверяем, есть ли уже значение по этому ключу
	if val, ok := c.buckets[hash].data[key]; ok {
		// если значение есть
		// перемещаем значение в начало очереди
		c.queue.MoveToFront(val)
		// обновляем значение
		val.Value.(*cacheItem).value = value
		return true
	}

	// если достигли максимального размера кэша, то удаляем последний элемент из очереди
	if c.queue.Len() == c.capacity {
		c.deleteOldest()
	}

	// устанавливаем срок жизни элемента
	ttl := time.Now().Add(c.ttl)

	// создаём элемент кэша
	item := &cacheItem{
		key:   key,
		value: value,
		ttl:   &ttl,
	}

	// записываем элемент в начало очереди
	el := c.queue.PushFront(item)
	c.buckets[hash].data[item.key] = el

	return true
}

// получаем значение из кэша
func (c *cache) Get(ctx context.Context, key string) interface{} {
	// определяем индекс бакета
	hash := hashValue(key, len(c.buckets))

	// блокировка
	c.buckets[hash].mu.Lock()
	defer c.buckets[hash].mu.Unlock()

	// пытаемся прочитать значение по ключу
	if el, ok := c.buckets[hash].data[key]; ok {
		// если прочитали успешно
		// перемещаем значение в начало очереди
		c.queue.MoveToFront(el)
		// возвращаем значение
		return el.Value.(*cacheItem).value
	}

	// не удалось найти значение в кэше
	return nil
}

// удаляем последний элемент из очереди
func (c *cache) deleteOldest() {
	if el := c.queue.Back(); el != nil {
		c.deleteCacheItem(el)
	}
}

func (c *cache) cleanCache() {
	for {
		<-c.ticker.C

		// итерируемся по всем бакетам
		for _, bucket := range c.buckets {
			// блокировка
			bucket.mu.Lock()
			// итерируемся по всем элементам
			for _, val := range bucket.data {
				// приводим элемент к нашему типу
				el := val.Value.(*cacheItem)
				if el.ttl != nil {
					// если срок жизни элемента истёк, то удаляем
					if time.Now().After(*el.ttl) {
						c.deleteCacheItem(val)
					}
				}
			}
			bucket.mu.Unlock()
		}

	}
}

// удаляем элемент кэша
func (c *cache) deleteCacheItem(el *list.Element) {
	// удаляем из очереди
	item := c.queue.Remove(el).(*cacheItem)
	// определяем индекс бакета
	hash := hashValue(item.key, len(c.buckets))
	// удалем из map
	delete(c.buckets[hash].data, item.key)
}
