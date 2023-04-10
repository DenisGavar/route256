package cache

import (
	"container/list"
	"sync"
	"time"
)

// элемент кэша
type cacheItem struct {
	key   string
	value interface{}

	ttl *time.Time
}

type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) interface{}
}

// кэш
type cache struct {
	// ёмкость кэша
	capacity int
	// ключ - значение
	data map[string]*list.Element
	//очередь для LRU
	queue *list.List
	// тикер для очистки кэша по ttl
	ticker *time.Ticker
	// ttl по умолчанию
	ttl time.Duration
	// для thread safe
	mu *sync.Mutex
}

// конструктор
func NewCache(capacity int, ttl time.Duration) Cache {
	cache := &cache{
		capacity: capacity,
		data:     make(map[string]*list.Element),
		queue:    list.New(),
		ticker:   time.NewTicker(time.Second * 10),
		ttl:      ttl,
		mu:       &sync.Mutex{},
	}

	// запускаем горутину для очистки кэша по ttl
	go cache.cleanCache()

	return cache
}

// записываем значение в кэш
func (c *cache) Set(key string, value interface{}) bool {
	// блокировка
	c.mu.Lock()
	defer c.mu.Unlock()

	// проверяем, есть ли уже значение по этому ключу
	if val, ok := c.data[key]; ok {
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
	c.data[item.key] = el

	return true
}

// получаем значение из кэша
func (c *cache) Get(key string) interface{} {
	// блокировка
	c.mu.Lock()
	defer c.mu.Unlock()

	// пытаемся прочитать значение по ключу
	if el, ok := c.data[key]; ok {
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
		// блокировка
		c.mu.Lock()
		// итерируемся по всем элементам
		for _, val := range c.data {
			// приводим элемент к нашему типу
			el := val.Value.(*cacheItem)
			if el.ttl != nil {
				// если срок жизни элемента истёк, то удаляем
				if time.Now().After(*el.ttl) {
					c.deleteCacheItem(val)
				}
			}
		}
		c.mu.Unlock()
	}
}

// удаляем элемент кэша
func (c *cache) deleteCacheItem(el *list.Element) {
	// удаляем из очереди
	item := c.queue.Remove(el).(*cacheItem)
	// удалем из map
	delete(c.data, item.key)
}
