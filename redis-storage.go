package quota

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	REDIS_KEY_PREFIX     = "holiday"      // The prefix for all keys.
	REDIS_TIMESTAMPS_KEY = "%s:timestamp" // The key for the timestamps.
	REDIS_KEY_FORMAT     = "%s:%s:%s:%d"  // The format of the redis key. The first parameter is the prefix, the second is the symbol, the third is the interval and the fourth is the timestamp.
)

// RedisStorage is a storage implementation that uses Redis as a backend.
type RedisStorage struct {
	mutex    *sync.RWMutex
	ctx      context.Context
	client   *redis.Client
	symbol   string
	interval time.Duration
}

// NewRedisStorage creates a new RedisStorage instance.
func NewRedisStorage(symbol string, interval time.Duration, config map[string]interface{}, ctx context.Context) *RedisStorage {
	// Prepare the options.
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config["redisHost"], config["redisPort"]),
		Password: config["redisPassword"].(string),
		DB:       config["redisDB"].(int),
	}

	// Initialize the driver.
	return &RedisStorage{
		mutex:    &sync.RWMutex{},
		ctx:      ctx,
		client:   redis.NewClient(options),
		symbol:   symbol,
		interval: interval,
	}
}

// Implemenet the Storage interface for RedisStorage.
// All returns all the values.
func (s *RedisStorage) All() ([]*Candle, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get all the keys.
	timestamps, err := s.client.LRange(s.ctx, fmt.Sprintf(REDIS_TIMESTAMPS_KEY, REDIS_KEY_PREFIX), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the timestamps from redis: %s", err)
	}

	var candles []*Candle

	// Retrieve the values.
	for _, ts := range timestamps {
		timestamp, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the timestamp: %s", err)
		}
		key := fmt.Sprintf(REDIS_KEY_FORMAT, REDIS_KEY_PREFIX, s.symbol, s.interval, timestamp)
		value, err := s.client.Get(s.ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get the candle for %s: %v", key, err)
		}

		// Transform the result into a map.
		candle, err := DeserializeCandle(value)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize the candle for %s: %v", key, err)
		}
		candles = append(candles, candle)
	}

	return candles, nil
}

// Get retrieves the value for the given key.
func (s *RedisStorage) Get(openTime time.Time) (*Candle, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Retrieve the value.
	timestamp := openTime.Unix()
	key := fmt.Sprintf(REDIS_KEY_FORMAT, REDIS_KEY_PREFIX, s.symbol, s.interval, timestamp)
	value, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Transform the result into a map.
	candle, err := DeserializeCandle(value)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize the candle for %s: %v", key, err)
	}

	return candle, nil
}

// Put stores the value for the given key.
func (s *RedisStorage) Put(c ...*Candle) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Prepare the data.
	for _, candle := range c {
		timestamp := candle.OpenTime.Unix()
		key := fmt.Sprintf(REDIS_KEY_FORMAT, REDIS_KEY_PREFIX, s.symbol, s.interval, timestamp)
		value, err := SerializeCandle(candle)
		if err != nil {
			return fmt.Errorf("failed to serialize the candle: %s", err)
		}

		// Store the value.
		err = s.client.Set(s.ctx, key, value, 0).Err()
		if err != nil {
			return fmt.Errorf("failed to store the value in redis: %s", err)
		}

		// Then store the key.
		err = s.client.LPush(s.ctx, fmt.Sprintf(REDIS_TIMESTAMPS_KEY, REDIS_KEY_PREFIX), timestamp).Err()
		if err != nil {
			return fmt.Errorf("failed to store the key in redis: %s", err)
		}
	}

	return nil
}

// Update updates the value for the given key.
func (s *RedisStorage) Update(candle *Candle) error {
	return s.Put(candle)
}

// Delete removes the value for the given key.
func (s *RedisStorage) Delete(c *Candle) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete the value.
	timestamp := c.OpenTime.Unix()
	key := fmt.Sprintf(REDIS_KEY_FORMAT, REDIS_KEY_PREFIX, s.symbol, s.interval, timestamp)
	err := s.client.Del(s.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete the value from redis: %s", err)
	}

	// Delete the key.
	err = s.client.LRem(s.ctx, fmt.Sprintf(REDIS_TIMESTAMPS_KEY, REDIS_KEY_PREFIX), 0, timestamp).Err()
	if err != nil {
		return fmt.Errorf("failed to delete the key from redis: %s", err)
	}

	return nil
}

// Close closes the RedisStorage instance.
func (s *RedisStorage) Close() error {
	return s.client.Close()
}

// PersistOlds will store the old candles into a persistance storage and remove them from the quota.
func (s *RedisStorage) PersistOlds(persist Storage, size int) error {
	// redis is a persistent storage, so we don't need to persist olds.
	return nil
}
