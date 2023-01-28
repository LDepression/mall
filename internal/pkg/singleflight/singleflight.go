package singleflight

import "sync"

/*
缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。
缓存击穿：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。
缓存穿透：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。
*/

// call 代表正在进行中，或已经结束的请求。使用 sync.WaitGroup 锁避免重入。
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// group 是 singleflight 的主数据结构，管理不同 key 的请求(call)。
type group struct {
	mu sync.Mutex
	m  map[interface{}]*call
}

var Group *group
var once sync.Once

func init() {
	once.Do(func() {
		Group = new(group)
	})
}

// Do 保证key所对应的fn函数同一时刻只会执行一次
func (g *group) Do(key interface{}, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	// 如果发现有函数正在运行,则等待其运行并返回其返回值
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 唯一一个运行的函数,添加到map中,并且设置锁
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()
	// 完成任务后从map中删除
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}

// 并发协程之间不需要消息传递，非常适合 sync.WaitGroup
