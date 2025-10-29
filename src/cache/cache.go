package cache

import (
	"sync"
	"time"
)

var Catch = NewExpiringMap()

// 定义缓存结构体：包含map（存储键值+过期时间）、读写锁（线程安全）
type ExpiringMap struct {
	data  map[interface{}]cacheItem // key:任意类型，value:包含值和过期时间的结构体
	mutex sync.RWMutex              // 读写锁：读共享、写排他
}

// cacheItem：存储单个键的value和过期时间戳（单位：纳秒，time.Time.UnixNano()）
type cacheItem struct {
	value  interface{} // 实际存储的值（支持任意类型）
	expire int64       // 过期时间戳（纳秒），0表示永不过期
}

// NewExpiringMap：创建一个新的带过期时间的缓存Map
func NewExpiringMap() *ExpiringMap {
	return &ExpiringMap{
		data: make(map[interface{}]cacheItem), // 初始化map
	}
}

// Set：添加/更新缓存键
// 参数：
// - key：缓存键（任意类型，如string、int）
// - value：缓存值（任意类型）
// - ttl：过期时间（time.Duration，如time.Hour表示1小时过期；ttl<=0表示永不过期）
func (em *ExpiringMap) Set(key, value interface{}, ttl time.Duration) {
	em.mutex.Lock()         // 写操作加排他锁，防止并发修改
	defer em.mutex.Unlock() // 函数结束后自动释放锁

	// 计算过期时间戳：若ttl<=0，expire设为0（永不过期）
	var expire int64
	if ttl > 0 {
		expire = time.Now().Add(ttl).UnixNano() // 当前时间+ttl，转为纳秒戳
	}

	// 存入map
	em.data[key] = cacheItem{
		value:  value,
		expire: expire,
	}
}

// Get：获取缓存键的值
// 返回：
// - value：缓存值（若键不存在或过期，返回nil）
// - ok：true表示键存在且未过期，false表示不存在或过期
func (em *ExpiringMap) Get(key interface{}) (value interface{}, ok bool) {
	em.mutex.RLock()         // 读操作加共享锁，支持多协程同时读
	defer em.mutex.RUnlock() // 函数结束后释放锁

	// 1. 检查键是否存在
	item, exists := em.data[key]
	if !exists {
		return nil, false // 键不存在
	}

	// 2. 检查是否过期（expire=0表示永不过期）
	now := time.Now().UnixNano()
	if item.expire != 0 && now > item.expire {
		return nil, false // 键已过期
	}

	go em.CleanExpired()

	// 3. 未过期，返回值
	return item.value, true
}

// Delete：删除指定缓存键（无论是否过期）
func (em *ExpiringMap) Delete(key interface{}) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	delete(em.data, key) // 直接从map删除键
}

// CleanExpired：主动清理所有过期键（可定期调用，如每小时调用一次）
func (em *ExpiringMap) CleanExpired() int {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	now := time.Now().UnixNano()
	delCount := 0 // 记录删除的过期键数量

	// 遍历map，删除所有过期键
	for key, item := range em.data {
		if item.expire != 0 && now > item.expire {
			delete(em.data, key)
			delCount++
		}
	}

	return delCount // 返回本次清理的过期键数量
}

// Len：获取当前缓存中“未过期键”的数量（需遍历检查过期）
func (em *ExpiringMap) Len() int {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	now := time.Now().UnixNano()
	count := 0

	for _, item := range em.data {
		if item.expire == 0 || now <= item.expire {
			count++
		}
	}

	return count
}
