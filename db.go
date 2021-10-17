// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package timezoneLookup

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/evanoberholster/timezoneLookup/pb"
	json "github.com/goccy/go-json"
	"github.com/klauspost/compress/snappy"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
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

func (dst *PolygonIndex) FromPB(src *pb.PolygonIndex) {
	dst.Id, dst.Tzid = src.Id, src.Tzid
	dst.Max.FromPB(src.Max)
	dst.Min.FromPB(src.Min)
}
func (src *PolygonIndex) ToPB(dst *pb.PolygonIndex) {
	dst.Reset()
	dst.Id, dst.Tzid = src.Id, src.Tzid
	dst.Max = src.Max.ToPB(dst.Max)
	dst.Min = src.Min.ToPB(dst.Min)
}

func BoltdbStorage(snappy bool, filename string, encoding encoding) TimezoneInterface {
	filename += "." + encoding.String()
	if snappy {
		filename += ".snap"
	}
	filename += ".db"
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
	switch e {
	case EncMsgPack:
		return "msgpack"
	case EncJSON:
		return "json"
	case EncProtobuf:
		return "protobuf"
	default:
		return "unknown"
	}
}
func EncodingFromString(s string) (encoding, error) {
	switch s {
	case "msgpack":
		return EncMsgPack, nil
	case "json":
		return EncJSON, nil
	case "protobuf":
		return EncProtobuf, nil
	default:
		return EncUnknown, fmt.Errorf("unknown encoding %q (neither msgpack, nor json)", s)
	}
}

var (
	EncUnknown  = encoding{}
	EncMsgPack  = encoding{1}
	EncJSON     = encoding{2}
	EncProtobuf = encoding{3}
)

func (s *Store) LoadTimezones() error {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return errors.New(errNotExistDatabase)
	}
	err := s.OpenDB(s.filename)
	if err != nil {
		return err
	}
	var pbIndex pb.PolygonIndex
	var U func(index *PolygonIndex, v []byte) error
	switch s.encoding {
	case EncMsgPack:
		U = func(index *PolygonIndex, v []byte) error {
			return msgpack.Unmarshal(v, index)
		}
	case EncJSON:
		U = func(index *PolygonIndex, v []byte) error {
			return json.Unmarshal(v, index)
		}
	case EncProtobuf:
		U = func(index *PolygonIndex, v []byte) error {
			if err := proto.Unmarshal(v, &pbIndex); err != nil {
				return err
			}
			index.FromPB(&pbIndex)
			return nil
		}
	}
	// Load polygon indexes
	return s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Index"))

		return b.ForEach(func(k, v []byte) error {
			var index PolygonIndex
			if err := U(&index, v); err != nil {
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
	var bufPolygon, bufIndex []byte
	var E func(polygon Polygon, index PolygonIndex) ([]byte, []byte, error)
	switch s.encoding {
	case EncMsgPack:
		pBuf, iBuf := bytes.NewBuffer(bufPolygon), bytes.NewBuffer(bufIndex)
		eP := msgpack.NewEncoder(pBuf)
		eI := msgpack.NewEncoder(iBuf)
		E = func(polygon Polygon, index PolygonIndex) ([]byte, []byte, error) {
			pBuf.Reset()
			if err := eP.Encode(polygon); err != nil {
				return nil, nil, err
			}
			// Marshal Polygon Index
			iBuf.Reset()
			err := eI.Encode(index)
			return pBuf.Bytes(), iBuf.Bytes(), err
		}
	case EncJSON:
		pBuf, iBuf := bytes.NewBuffer(bufPolygon), bytes.NewBuffer(bufIndex)
		eP := json.NewEncoder(pBuf)
		eI := json.NewEncoder(iBuf)
		E = func(polygon Polygon, index PolygonIndex) ([]byte, []byte, error) {
			pBuf.Reset()
			if err := eP.Encode(polygon); err != nil {
				return nil, nil, err
			}
			iBuf.Reset()
			err := eI.Encode(index)
			return pBuf.Bytes(), iBuf.Bytes(), err
		}
	case EncProtobuf:
		var pbPoly pb.Polygon
		var pbIndex pb.PolygonIndex
		var mo proto.MarshalOptions
		E = func(polygon Polygon, index PolygonIndex) ([]byte, []byte, error) {
			polygon.ToPB(&pbPoly)
			bufPolygon, err := mo.MarshalAppend(bufPolygon[:0], &pbPoly)
			if err != nil {
				return nil, nil, err
			}
			index.ToPB(&pbIndex)
			bufIndex, err := mo.MarshalAppend(bufIndex[:0], &pbIndex)
			return bufPolygon, bufIndex, err
		}
	}
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
			bufPolygon, bufIndex, err := E(polygon, index)
			if err != nil {
				return err
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
	var pbPoly pb.Polygon
	var U func(polygon *Polygon, v []byte) error
	switch s.encoding {
	case EncMsgPack:
		U = func(polygon *Polygon, v []byte) error {
			return msgpack.Unmarshal(v, polygon)
		}
	case EncJSON:
		U = func(polygon *Polygon, v []byte) error {
			return json.Unmarshal(v, polygon)
		}
	case EncProtobuf:
		U = func(polygon *Polygon, v []byte) error {
			if err := proto.Unmarshal(v, &pbPoly); err != nil {
				return err
			}
			polygon.FromPB(&pbPoly)
			return nil
		}
	}
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
		return U(&polygon, v)
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
