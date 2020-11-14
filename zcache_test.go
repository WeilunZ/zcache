package zcache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, err := f.Get("key"); err != nil || !reflect.DeepEqual(v, expect) {
		t.Fatalf("callback failed")
	}
}

