package config

import (
	"encoding/json"
	"os"

	"golang.org/x/exp/slog"
)

type kv struct {
	file string
	val  map[int]bool
}

func newKV(file string) *kv {
	instance := &kv{file: file}

	data, err := os.ReadFile(file)
	if err != nil {
		slog.Error("[Config] Failed to create kv store", "err", err)
		return instance
	}

	json.Unmarshal(data, &instance.val)
	return instance
}

func (instance *kv) getBool(id int, dft bool) bool {
	if val, ok := instance.val[id]; ok {
		return val
	}
	return dft
}

func (instance *kv) setBool(id int, enabled bool) {
	instance.val[id] = enabled

	data, err := json.Marshal(instance.val)
	if err != nil {
		slog.Error("[Config] Failed to persist kv store", "err", err)
		return
	}

	os.WriteFile(instance.file, data, 0666)
}
