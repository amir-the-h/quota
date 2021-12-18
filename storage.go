package quota

import (
	"time"
)

// Storage is an interface for storing and retrieving candles.
type Storage interface {
	// All returns all the candle in the storage.
	All(symbol string, interval time.Duration) (*Quota, error)

	// Get retrieves candle from the storage.
	Get(openTime time.Time) (*Candle, error)

	// GetByIndex retrieves candle from the storage by index.
	GetByIndex(index int) (*Candle, error)

	// Put stores the given candle in the storage.
	Put(candle ...*Candle) error

	// Update updates candle in the storage.
	Update(candle ...*Candle) error

	// Delete removes the candle from the storage.
	Delete(key string) error

	// Close closes the storage.
	Close() error

	// PersistOlds will store the old candles into a persistance storage and remove them from the quota.
	PersistOlds(persist Storage, size int) error
}
