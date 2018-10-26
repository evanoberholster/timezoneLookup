package timezoneLookup
import (
	"log"
	"errors"
	"encoding/binary"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"github.com/golang/snappy"
)

type Store struct { 	// Database struct
	db 			*bolt.DB
	pIndex 		[]PolygonIndex
	filename 	string
	snappy 		bool
}

type PolygonIndex struct {
	Id 			uint64 		`json:"-"`
	Tzid 		string		`json:"tzid"`
	Max		    Coord 		`json:"max"`
	Min 	    Coord 		`json:"min"`
}

func BoltdbStorage(snappy bool, filename string) *Store {
	if snappy {
		filename = filename + ".snap.db"
	} else {
		filename = filename + ".db"
	}
	return &Store{
		filename: filename,
		pIndex: []PolygonIndex{},
		snappy: snappy,
	}
}

func (s *Store)Close() {
	defer s.db.Close()
	PrintMemUsage() 
}

func (s *Store)LoadTimezones() (error) {
	err := s.OpenDB(s.filename)
	if err != nil {
		return err
	}
	// Load indexes 
	return s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Index"))
		
		b.ForEach(func(k, v []byte) error {
			var index PolygonIndex
			err := json.Unmarshal(v, &index)
			if err != nil {
				log.Println(err)
			}
			index.Id = binary.BigEndian.Uint64(k)
			s.pIndex = append(s.pIndex, index)
			return nil
		})
		return nil
	})
}

func (s *Store)Query(q Coord) (string, error) {
	for _, i := range s.pIndex {
		if i.Min.X < q.X && i.Min.Y < q.Y && i.Max.X > q.X && i.Max.Y > q.Y {
			p, err := s.loadPolygon(i.Id)
			if err != nil {
				return "Error", err
			} 
			if p.contains(q) {
				return i.Tzid, nil
			}
		}
	}
	return "Error", errors.New("Timezone not found")
}

func (s *Store)CreateTimezones(jsonFilename string) (error)  {
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
		s.InsertPolygons(tz)
	}
	return nil
}

func (s *Store)createBuckets() (error) {
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

func (s *Store)InsertPolygons(tz Timezone) {
	for _, polygon := range tz.Polygons {
		s.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Polygon"))
			i := tx.Bucket([]byte("Index"))

			// Get ID number autoIncrement
			id, _ := b.NextSequence()
        	intId := int(id)

        	// Create Polygon Index
        	index := PolygonIndex{
        		Tzid: tz.Tzid,
				Max: polygon.Max,
				Min: polygon.Min,
        	}

        	// UnMarshal Polygon Index
        	bufPolygon, err := json.Marshal(polygon)
		    if err != nil {
		        return err
		    }
		    // UnMarshal Polygon
        	bufIndex, err := json.Marshal(index)
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
		})
	}
}

func (s *Store)loadPolygon(id uint64) (Polygon, error) {
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
		return json.Unmarshal(v, &polygon)
	})
	return polygon, err
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}


func (s *Store)OpenDB(path string) (error) {
	var err error
	s.db, err = bolt.Open(path, 0666, nil)
	if err != nil {
	  return err
	}
	return nil
}
