package webber

import (
	"fmt"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const (
	defaultDBName = "default"
)

// Helper type for a map[string]string (will be a map[string]any once more types are supported)
type M map[string]string

type Collection struct {
	*bbolt.Bucket
}

type Webber struct {
	db *bbolt.DB
}

func New() (*Webber, error) {
	dbName := fmt.Sprintf("%s.webber", defaultDBName)
	db, err := bbolt.Open(dbName, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &Webber{
		db: db,
	}, nil
}

func (w *Webber) CreateCollectionIfNotExists(name string) (*Collection, error) {
	tx, err := w.db.Begin(true)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	coll := Collection{}
	bucket, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	coll.Bucket = bucket

	if err != nil {
		return nil, err
	}
	return &coll, nil
}

func (w *Webber) Insert(collName string, data M) (uuid.UUID, error) {
	id := uuid.New()

	tx, err := w.db.Begin(true)
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	collBucket, err := tx.CreateBucketIfNotExists([]byte(collName))
	if err != nil {
		return id, err
	}

	recordBucket, err := collBucket.CreateBucket([]byte(id.String()))
	if err != nil {
		return id, err
	}

	for k, v := range data {
		if err := recordBucket.Put([]byte(k), []byte(v)); err != nil {
			return id, err
		}
	}

	if err := recordBucket.Put([]byte("id"), []byte(id.String())); err != nil {
		return id, err
	}

	return id, err
}

func (w *Webber) Select(coll, k string, query any) error {
	return nil
}
