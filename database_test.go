package jinx_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hedarikun/jinx"
)

func TestGet(t *testing.T) {
	db := jinx.New()
	db.Set("foo", "bar")

	val1 := db.Get("test")
	val2 := db.Get("foo").(string)

	if val1 != nil {
		t.Errorf("expected nil, got %v", val1)
	}
	if val2 != "bar" {
		t.Errorf("expected bar, got %v", val2)
	}
}

func TestTransactionCorrect(t *testing.T) {
	db := jinx.New()
	db.Set("foo", "bar")

	err := db.HandleTransaction(func(tx *jinx.JinxTransaction) error {
		tx.Set("foo", "baz")
		tx.Set("test", "test2")
		return nil
	})

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	val1 := db.Get("foo").(string)
	val2, ok := db.Get("test").(string)

	if !ok {
		t.Errorf("expected true, got %v", ok)
	}

	if val1 != "baz" {
		t.Errorf("expected baz, got %v", val1)
	}
	if val2 != "test2" {
		t.Errorf("expected test2, got %v", val2)
	}
}

func TestTransactionIncorrect(t *testing.T) {
	db := jinx.New()
	db.Set("foo", "bar")

	err := db.HandleTransaction(func(tx *jinx.JinxTransaction) error {
		tx.Set("foo", "baz")
		tx.Set("test", "test2")
		return errors.New("test error")
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	val1 := db.Get("foo").(string)
	val2, ok := db.Get("test").(string)

	if ok {
		t.Errorf("expected false, got %v", ok)
	}

	if val1 != "bar" {
		t.Errorf("expected bar, got %v", val1)
	}
	if val2 != "" {
		t.Errorf("expected empty string, got %v", val2)
	}
}

func TestSetExpire(t *testing.T) {
	db := jinx.New()
	db.SetExpire("foo", "bar", 3)

	time.Sleep(3 * time.Second)

	val := db.Get("foo")

	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}

	db.SetExpire("foo", "bar", 2)

	val = db.Get("foo")

	if val == nil {
		t.Errorf("expected bar, got %v", val)
	}
}
