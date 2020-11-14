package lru

import "container/list"

type Cache struct {
	maxBytes uint64 //允许使用的最大内存
	nbytes uint64 //当前已使用的内存
	cache map[string]*list.Element
	dlist *list.List
	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int
}

// New creates a lru cache
func New(maxBytes uint64, onEvicted func(string, Value)) *Cache{
	return &Cache{
		maxBytes:  maxBytes,
		cache:     make(map[string]*list.Element),
		dlist:     list.New(),
		OnEvicted: onEvicted,
	}
}

// Get gets the value of the given key
// return nil, false if not exists
func (c *Cache) Get(key string) (Value, bool){
	if element, ok := c.cache[key]; ok {
		c.dlist.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// RemoveOldest removes the least recently used entry
func (c *Cache) RemoveOldest(){
	element := c.dlist.Back()
	if element != nil{
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.dlist.Remove(element)
		c.nbytes -= uint64(len(kv.key)) + uint64(kv.value.Len())
		if c.OnEvicted != nil{
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Put [key, value] pair
func (c *Cache) Put(key string, val Value){
	if v, ok := c.cache[key]; !ok{
		ele := c.dlist.PushFront(&entry{
			key:   key,
			value: val,
		})
		c.cache[key] = ele
		c.nbytes += uint64(len(key)) + uint64(val.Len())
	}else{
		c.dlist.MoveToFront(v)
		kv := v.Value.(*entry)
		c.nbytes += uint64(val.Len()) - uint64(kv.value.Len())
		kv.value = val
	}
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.dlist.Len()
}












