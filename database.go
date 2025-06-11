package jinx

import (
	"sync"
	"time"
)

type JinxDatabase struct {
	entries map[string]databaseEntry
	mu      sync.RWMutex
}

func New() *JinxDatabase {
	return &JinxDatabase{
		entries: make(map[string]databaseEntry),
	}
}

func (jd *JinxDatabase) Set(key string, value interface{}) error {
	jd.mu.Lock()
	defer jd.mu.Unlock()

	var rValue interface{}

	switch v := value.(type) {
	case string:
		rValue = v
	case []byte:
		rValue = string(v)
	default:
		rValue = value
	}

	databaseEntry := databaseEntry{
		Value:          rValue,
		expirationDate: 0,
	}
	jd.entries[key] = databaseEntry
	return nil
}

// expiration is time in seconds
func (jd *JinxDatabase) SetExpire(key string, value interface{}, expire int) error {
	jd.mu.Lock()
	defer jd.mu.Unlock()

	var rValue interface{}

	switch v := value.(type) {
	case string:
		rValue = v
	case []byte:
		rValue = string(v)
	default:
		rValue = value
	}

	databaseEntry := databaseEntry{
		Value:          rValue,
		expirationDate: time.Now().Add(time.Duration(expire) * time.Second).Unix(),
	}
	jd.entries[key] = databaseEntry
	return nil
}

func (jd *JinxDatabase) SetMap(key string, m map[string]string) error {
	jd.mu.Lock()
	defer jd.mu.Unlock()
	jd.entries[key] = databaseEntry{
		Value:          m,
		expirationDate: 0,
	}
	return nil
}

func (jd *JinxDatabase) Get(key string) interface{} {
	jd.mu.RLock()
	defer jd.mu.RUnlock()

	element, ok := jd.entries[key]
	if !ok {
		return nil
	}

	if element.expirationDate != 0 && element.expirationDate <= time.Now().Unix() {
		delete(jd.entries, key)
		return nil
	}

	return jd.entries[key].Value
}

func (jd *JinxDatabase) Delete(key string) {
	jd.mu.Lock()
	defer jd.mu.Unlock()
	delete(jd.entries, key)
}

func (jd *JinxDatabase) Exists(key string) bool {
	jd.mu.RLock()
	defer jd.mu.RUnlock()
	_, ok := jd.entries[key]
	return ok
}

func (jd *JinxDatabase) Merge(otherDB *JinxDatabase) error {
	otherDB.mu.RLock()
	defer otherDB.mu.RUnlock()

	jd.mu.Lock()
	defer jd.mu.Unlock()

	for key, value := range otherDB.entries {
		jd.entries[key] = value
	}
	return nil
}

func (jd *JinxDatabase) KeyCount() int {
	jd.mu.RLock()
	defer jd.mu.RUnlock()
	return len(jd.entries)
}

func (jd *JinxDatabase) HandleTransaction(handler TransactionHandler) error {
	localDB := New()
	tx := JinxTransaction{
		localDB:  localDB,
		parentDB: jd,
	}

	err := handler(&tx)
	if err != nil {
		return err
	}

	jd.Merge(localDB)
	return nil
}
