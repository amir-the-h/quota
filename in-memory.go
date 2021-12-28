package quota

import (
	"fmt"
	"sync"
	"time"
)

// InMemoryStorage is an in-memory implementation of the Storage interface.
type InMemoryStorage struct {
	mutex    *sync.RWMutex
	symbol   string
	interval time.Duration
	Q        *Quota
}

// NewInMemoryStorage creates a new InMemoryStorage instance.
func NewInMemoryStorage(symbol string, interval time.Duration) *InMemoryStorage {
	return &InMemoryStorage{
		mutex:    &sync.RWMutex{},
		symbol:   symbol,
		interval: interval,
		Q:        &Quota{},
	}
}

// All returns all the data in the storage.
func (s *InMemoryStorage) All() (*Quota, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Q, nil
}

// Get returns the value for the given key.
func (s *InMemoryStorage) Get(openTime time.Time) (*Candle, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	candle, _ := s.Q.Find(openTime.Unix())
	return candle, nil
}

// GetByIndex retrieves candle from the storage by index.
func (s *InMemoryStorage) GetByIndex(index int) (*Candle, error) {
	return (*s.Q)[index], nil
}

// Put stores the given value for the given key.
func (s *InMemoryStorage) Put(c ...*Candle) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, candle := range c {
		*s.Q = append(*s.Q, candle)
	}
	return nil
}

// Update updates the value for the given key.
func (s *InMemoryStorage) Update(c ...*Candle) error {
	return s.Put(c...)
}

// Delete deletes the value for the given key.
func (s *InMemoryStorage) Delete(c *Candle) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, index := s.Q.Find(c.OpenTime.Unix())
	if index == -1 {
		return fmt.Errorf("candle not found")
	}
	*s.Q = append((*s.Q)[:index], (*s.Q)[index+1:]...)
	return nil
}

// Close closes the storage.
func (s *InMemoryStorage) Close() error {
	*s.Q = (*s.Q)[:0]
	return nil
}

// PersistOlds will store the old candles into a persistance storage and remove them from the quota.
func (s *InMemoryStorage) PersistOlds(persist Storage, size int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(*s.Q) <= size {
		return fmt.Errorf("not enough candles to persist")
	}
	candles := (*s.Q)[:size]
	*s.Q = (*s.Q)[size:]
	return persist.Put(candles...)
}
