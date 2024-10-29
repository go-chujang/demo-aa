package config

import (
	"fmt"
	"math/big"
	"strconv"
	"time"
)

func get(k key) (string, error) {
	v, exist := confHolder[k]
	if !exist {
		return "", ErrEmptyString
	}
	return v, nil
}

func Get(k key) string {
	v, _ := get(k)
	return v
}

func GetInt64(k key) (int64, error) {
	v, err := get(k)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 64)
}

func GetInt(k key) (int, error) {
	v, err := get(k)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(v)
}

func GetUint64(k key) (uint64, error) {
	v, err := get(k)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(v, 10, 64)
}

func GetDuration(k key) (time.Duration, error) {
	v, err := get(k)
	if err != nil {
		return 0, nil
	}
	return time.ParseDuration(v)
}

func GetDurationSeconds(k key) (int64, error) {
	t, err := GetDuration(k)
	if err != nil {
		return 0, nil
	}
	return int64(t.Seconds()), nil
}

func GetBig(k key) (*big.Int, error) {
	v, err := get(k)
	if err != nil {
		return nil, err
	}
	b, ok := new(big.Int).SetString(v, 10)
	if !ok {
		return nil, fmt.Errorf("failed parse to big.Int key: %s, value: %v", k, v)
	}
	return b, nil
}
