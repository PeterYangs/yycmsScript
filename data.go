package yycmsScript

import "sync"

type Data struct {
	store sync.Map
}

func NewData() *Data {

	return &Data{store: sync.Map{}}
}

func (d *Data) Set(key string, value string) {

	d.store.Store(key, value)

}

func (d *Data) Get(key string, defaultValue string) string {

	v, ok := d.store.Load(key)

	if !ok {

		return defaultValue
	}

	return v.(string)
}
