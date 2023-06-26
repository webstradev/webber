package webbr

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
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

const (
	defaultDBName = "default"
)

// Helper type for a map[string]string (will be a map[string]any once more types are supported)
type M map[string]any

type Filter struct {
	EQ    M
	Limit int
	Sort  string
}

type Collection struct {
	*bbolt.Bucket
}

type webbr struct {
	db *bbolt.DB
}

func New() (*webbr, error) {
	dbName := fmt.Sprintf("%s.webbr", defaultDBName)
	db, err := bbolt.Open(dbName, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &webbr{
		db: db,
	}, nil
}

func (w *webbr) CreateCollectionIfNotExists(name string) (*Collection, error) {
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

func (w *webbr) Insert(collName string, data M) (uuid.UUID, error) {
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
		typeInfo, err := getValueTypeInfo(v)
		if err != nil {
			return id, err
		}
		fmt.Printf("%+v\n", typeInfo)
		if err := recordBucket.Put([]byte(k), typeInfo.underlying); err != nil {
			return id, err
		}
	}

	if err := recordBucket.Put([]byte("id"), []byte(id.String())); err != nil {
		return id, err
	}

	return id, tx.Commit()
}

func (w *webbr) Find(coll string, filter Filter) ([]M, error) {
	tx, err := w.db.Begin(false)
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
		if v == nil {
			entryBucket := bucket.Bucket(k)
			if entryBucket == nil {
				return fmt.Errorf("entry found without field data")
			}

			data := M{}
			entryBucket.ForEach(func(k, v []byte) error {
				data[string(k)] = string(v)
				return nil
			})
			include := true
			if filter.EQ != nil {
				include = false
				for filterKey, filterValue := range filter.EQ {
					if value, ok := data[filterKey]; ok {
						if value == filterValue {
							include = true
						}
					}
				}
			}
			if include {
				results = append(results, data)
			}
		}
		return nil
	})

	return results, nil
}

type ValueTypeInfo struct {
	valueType  ValueType
	underlying []byte
}

func getValueTypeInfo(value any) (ValueTypeInfo, error) {
	switch it := value.(type) {
	case string:
		return ValueTypeInfo{
			valueType:  ValueTypeString,
			underlying: []byte(it),
		}, nil
	case int:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(it))
		return ValueTypeInfo{
			valueType:  ValueTypeInt,
			underlying: b,
		}, nil
	case float64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(it))
		return ValueTypeInfo{
			valueType:  ValueTypeFloat,
			underlying: b,
		}, nil
	case bool:
		var b []byte
		if it {
			b = []byte{0x01}
		} else {
			b = []byte{0x00}
		}
		return ValueTypeInfo{
			valueType:  ValueTypeBool,
			underlying: b,
		}, nil
	default:
		return ValueTypeInfo{
			valueType:  ValueTypeUnknown,
			underlying: []byte{},
		}, fmt.Errorf("unsupported type (%s)", reflect.TypeOf(it))
	}
}