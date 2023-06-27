package webbr

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

const (
	defaultDBName    = "default"
	defaultExtension = "webbr"
)

// Helper type for a map[string]any
type M map[string]any

type Filter struct {
	EQ     M
	Select []string
	Limit  int
	Sort   string
}

type Webbr struct {
	*Options
	db *bbolt.DB
}

func New(options ...OptFunc) (*Webbr, error) {
	opts := &Options{
		DBName:    defaultDBName,
		Extension: defaultExtension,
		Encoder:   JSONEncoder{},
		Decoder:   JSONDecoder{},
	}

	for _, fn := range options {
		fn(opts)
	}

	db, err := bbolt.Open(opts.GetDBName(), 0666, nil)
	if err != nil {
		return nil, err
	}

	return &Webbr{
		db:      db,
		Options: opts,
	}, nil
}

func (w *Webbr) DropDatabase(name string) error {
	return os.Remove(w.GetDBName())
}

func (w *Webbr) CreateCollection(name string) (*bbolt.Bucket, error) {
	tx, err := w.db.Begin(true)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (w *Webbr) Insert(collName string, data M) (uint64, error) {
	tx, err := w.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	collBucket, err := tx.CreateBucketIfNotExists([]byte(collName))
	if err != nil {
		return 0, err
	}

	id, err := collBucket.NextSequence()
	if err != nil {
		return 0, err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	if err := collBucket.Put(uint64ToBytes(id), b); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (w *Webbr) Update(coll string, filter Filter, data M) ([]M, error) {
	tx, err := w.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(coll))
	if bucket == nil {
		return nil, fmt.Errorf("collection (%s) not found", coll)
	}
	records, err := w.findFiltered(bucket, filter)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		for k, v := range data {
			if k == "id" {
				continue
			}
			if _, ok := record[k]; ok {
				record[k] = v
			}
		}

		b, err := w.Encoder.Encode(record)
		if err != nil {
			return nil, err
		}

		id := record["id"].(uint64)
		if err := bucket.Put(uint64ToBytes(id), b); err != nil {
			return nil, err
		}
	}
	return records, tx.Commit()
}

func (w *Webbr) Find(coll string, filter Filter) ([]M, error) {
	tx, err := w.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(coll))
	if bucket == nil {
		return nil, fmt.Errorf("collection (%s) not found", coll)
	}

	results, err := w.findFiltered(bucket, filter)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (w *Webbr) findFiltered(bucket *bbolt.Bucket, filter Filter) ([]M, error) {
	var records []M
	bucket.ForEach(func(k, v []byte) error {
		record := M{
			"id": uint64FromBytes(k),
		}
		if err := w.Decoder.Decode(v, &record); err != nil {
			return err
		}
		include := true
		// If there is an EQ Filter we will only grab records where that matches
		if filter.EQ != nil {
			include = false
			for filterKey, filterValue := range filter.EQ {
				if value, ok := record[filterKey]; ok {
					if filterValue == value {
						include = true
					}
				}
			}
		}
		if include {
			// If we are only selecting certain keys we filter them here
			if len(filter.Select) > 0 {
				data := M{}
				for _, k := range filter.Select {
					data[k] = record[k]
				}
				records = append(records, data)
			} else {
				records = append(records, record)
			}
		}
		return nil
	})
	return records, nil
}

// Helper functions
func uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

func uint64FromBytes(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}
