package store

import (
	"encoding/json"
	"fmt"
)

func marshalRecord(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode record: %w", err)
	}
	return encoded, nil
}

func unmarshalRecord(data []byte, target any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode record: %w", err)
	}
	return nil
}

func putRecord(bucketName []byte, key string, value any, put func([]byte, []byte) error) error {
	data, err := marshalRecord(value)
	if err != nil {
		return err
	}
	return put(entityKey(key), data)
}

func copyRecord(bucketName []byte, key []byte, read func([]byte) []byte) []byte {
	value := read(key)
	if value == nil {
		return nil
	}
	result := make([]byte, len(value))
	copy(result, value)
	return result
}
