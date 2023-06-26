package webber

import (
	"fmt"

	"go.etcd.io/bbolt"
)

const (
	defaultDBName = "default"
)

type Collection struct {
	bucket *bbolt.Bucket
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

func (w *Webber) CreateCollection(name string) (*Collection, error) {
	coll := Collection{}
	err := w.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		coll.bucket = bucket
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &coll, nil
}
