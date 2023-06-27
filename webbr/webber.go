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

type ValueType int

func (v ValueType) String() string {
	switch v {
	case ValueTypeUnknown:
		return "unknown"
	case ValueTypeString:
		return "string"
	case ValueTypeInt:
		return "integer"
	case ValueTypeBool:
		return "boolean"
	case ValueTypeFloat:
		return "float"
	}
	return "unknown"
}

const (
	ValueTypeUnknown = iota
	ValueTypeString
	ValueTypeInt
	ValueTypeBool
	ValueTypeFloat
)

// Helper type for a map[string]string (will be a map[string]any once more types are supported)
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

type OptFunc (func(opts *Options))

type Options struct {
	DBName    string
	Extension string
}

func (o Options) GetDBName() string {
	return fmt.Sprintf("%s.%s", o.DBName, o.Extension)
}

func WithDBName(name string) OptFunc {
	return func(opts *Options) {
		opts.DBName = name
	}
}

func WithExtension(ext string) OptFunc {
	return func(opts *Options) {
		opts.Extension = ext
	}
}

func New(options ...OptFunc) (*Webbr, error) {
	opts := &Options{
		DBName:    defaultDBName,
		Extension: defaultExtension,
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

	results := []M{}
	bucket.ForEach(func(k, v []byte) error {
		data := M{
			"id": uint64FromBytes(k),
		}
		if err := json.Unmarshal(v, &data); err != nil {
			return err
		}

		include := true
		if filter.EQ != nil {
			include = false
			for filterKey, filterValue := range filter.EQ {
				if value, ok := data[filterKey]; ok {
					if filterValue == value {
						include = true
					}
				}
			}
		}
		if include {
			results = append(results, data)
		}
		return nil
	})

	return results, nil
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
