package store

import (
	"encoding/binary"
	"errors"
	"path/filepath"

	"go.etcd.io/bbolt"

	"venue-reservation/internal/model"
)

var ErrNotFound = errors.New("record not found")

var bucketNames = [][]byte{
	[]byte("Venue"),
	[]byte("TimeSlot"),
	[]byte("Reservation"),
	[]byte("AuditEvent"),
	[]byte("CollaborationNote"),
	[]byte("ArchiveRecord"),
	[]byte("Metadata"),
}

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0o600, nil)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) update(fn func(*bbolt.Tx) error) error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(fn)
}

func (s *Store) view(fn func(*bbolt.Tx) error) error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.View(fn)
}

func (s *Store) bucket(tx *bbolt.Tx, name []byte) (*bbolt.Bucket, error) {
	bucket := tx.Bucket(name)
	if bucket == nil {
		return nil, errors.New("missing bucket")
	}
	return bucket, nil
}

func uintKey(value uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, value)
	return key
}

func (s *Store) nextSequence(tx *bbolt.Tx) (uint64, error) {
	bucket, err := s.bucket(tx, []byte("Metadata"))
	if err != nil {
		return 0, err
	}
	current := binary.BigEndian.Uint64(padded(bucket.Get([]byte("sequence"))))
	current++
	return current, bucket.Put([]byte("sequence"), uintKey(current))
}

func padded(value []byte) []byte {
	if len(value) == 8 {
		return value
	}
	return make([]byte, 8)
}

func entityKey(id string) []byte {
	return []byte(id)
}

func ensureVenue(venue model.Venue) error {
	return model.ValidateVenue(venue)
}
