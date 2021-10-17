// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package timezoneLookup

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/klauspost/compress/snappy"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
)

type Store struct { // Database struct
	db       *bolt.DB
	pIndex   []PolygonIndex
	filename string
	snappy   bool
	encoding encoding
}

type PolygonIndex struct {
	Id   uint64 `json:"-"`
	Tzid string `json:"tzid"`
	Max  Coord  `json:"max"`
	Min  Coord  `json:"min"`
}

func BoltdbStorage(snappy bool, filename string, encoding encoding) TimezoneInterface {
	if snappy {
		filename = filename + ".snap.db"
	} else {
		filename = filename + ".db"
	}
	return &Store{
		filename: filename,
		pIndex:   []PolygonIndex{},
		snappy:   snappy,
		encoding: encoding,
	}
}

func (s *Store) Close() {
	defer s.db.Close()
}

type encoding struct {
	Type uint8
}

func (e encoding) String() string {
	if e == EncMsgPack {
		return "msgpack"
	}
	return "json"
}
func EncodingFromString(s string) (encoding, error) {
	switch s {
	case "msgpack":
		return EncMsgPack, nil
	case "json":
		return EncJSON, nil
	default:
		return EncUnknown, fmt.Errorf("unknown encoding %q (neither msgpack, nor json)", s)
	}
}

var (
	EncUnknown = encoding{}
	EncMsgPack = encoding{1}
	EncJSON    = encoding{2}
)

func (s *Store) LoadTimezones() error {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return errors.New(errNotExistDatabase)
	}
	err := s.OpenDB(s.filename)
	if err != nil {
		return err
	}
	// Load polygon indexes
	return s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Index"))

		return b.ForEach(func(k, v []byte) error {
			var index PolygonIndex
			var err error
			if s.encoding == EncMsgPack {
				err = msgpack.Unmarshal(v, &index)
			} else {
				err = json.Unmarshal(v, &index)
			}
			if err != nil {
				return err
			}
			index.Id = binary.BigEndian.Uint64(k)
			s.pIndex = append(s.pIndex, index)
			return nil
		})
	})
}

func (s *Store) Query(q Coord) (string, error) {
	for _, i := range s.pIndex {
		if i.Min.Lat < q.Lat && i.Min.Lon < q.Lon && i.Max.Lat > q.Lat && i.Max.Lon > q.Lon {
			p, err := s.loadPolygon(i.Id)
			if err != nil {
				return i.Tzid, errors.New(errPolygonNotFound)
			}
			if p.contains(q) {
				return i.Tzid, nil
			}
		}
	}
	return "Error", errors.New(errTimezoneNotFound)
}

func (s *Store) CreateTimezones(jsonFilename string) error {
	err := checkFilesExist(jsonFilename, s.filename)
	if err != nil {
		return err
	}
	tzs, err := TimezonesFromGeoJSON(jsonFilename)
	if err != nil {
		return err
	}
	err = s.OpenDB(s.filename)
	if err != nil {
		return err
	}
	err = s.createBuckets()
	if err != nil {
		return err
	}
	for _, tz := range tzs {
		if err := s.InsertPolygons(tz); err != nil {
			return err
		}
	}
	return nil
}

func checkFilesExist(src string, dest string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New(errNotExistGeoJSON)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return errors.New(errExistDatabase)
	}
	return nil
}

func (s *Store) createBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Index"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte("Polygon"))
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) InsertPolygons(tz Timezone) error {
	for _, polygon := range tz.Polygons {
		if err := s.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Polygon"))
			i := tx.Bucket([]byte("Index"))

			// Get ID number autoIncrement
			id, _ := b.NextSequence()
			intId := int(id)

			// Create Polygon Index
			index := PolygonIndex{
				Tzid: tz.Tzid,
				Max:  polygon.Max,
				Min:  polygon.Min,
			}
			var bufPolygon, bufIndex []byte
			var err error

			if s.encoding == EncMsgPack {
				// Marshal Polygons
				bufPolygon, err = msgpack.Marshal(polygon)
				if err != nil {
					return err
				}
				// Marshal Polygon Index
				bufIndex, err = msgpack.Marshal(index)
				if err != nil {
					return err
				}
			} else {
				bufPolygon, err = json.Marshal(polygon)
				if err != nil {
					return err
				}
				bufIndex, err = json.Marshal(index)
				if err != nil {
					return err
				}
			}
			if s.snappy {
				bufPolygon = snappy.Encode(nil, bufPolygon)
			}
			// Write Polygon Index
			err = i.Put(itob(intId), bufIndex)
			if err != nil {
				return err
			}
			return b.Put(itob(intId), bufPolygon)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) loadPolygon(id uint64) (Polygon, error) {
	var polygon Polygon
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Polygon"))
		v := b.Get(itob(int(id)))
		if s.snappy {
			var err error
			v, err = snappy.Decode(nil, v)
			if err != nil {
				return err
			}
		}
		if s.encoding == EncMsgPack {
			return msgpack.Unmarshal(v, &polygon)
		} else {
			return json.Unmarshal(v, &polygon)
		}
	})
	return polygon, err
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (s *Store) OpenDB(path string) error {
	var err error
	s.db, err = bolt.Open(path, 0666, nil)
	if err != nil {
		return err
	}
	return nil
}
