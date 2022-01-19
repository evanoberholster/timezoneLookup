package timezoneLookup

import (
	"bufio"
	"encoding/binary"
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/evanoberholster/timezoneLookup/geo"
	"golang.org/x/sys/unix"
)

var (
	endian                 = binary.LittleEndian
	pageSize               = os.Getpagesize()
	ErrCoordinatesNotValid = errors.New("Latitude and/or Longitude are not valid")
)

type Timezonecache struct {
	data       []byte
	arr        []uint32
	name       []string
	rt         geo.RTree
	dataOffset uint32
	dataLength uint32
}

func (tzc *Timezonecache) AddTimezone(tz Timezone) {
	for _, p := range tz.Polygons {
		var offset uint32
		id := uint(len(tzc.arr)) // next id
		buf := p.ToByteSlice()
		tzc.data = append(tzc.data, buf...)
		if id == 0 {
			offset += tzc.dataOffset
		} else {
			offset += tzc.arr[id-1]
		}
		offset += uint32(len(buf))
		tzc.arr = append(tzc.arr, offset)
		tzc.name = append(tzc.name, tz.Name)
		tzc.rt.InsertPolygon(p, id)
	}
}

func (tzc *Timezonecache) buf(id uint) []byte {
	offset := uint32(0)
	if id == 0 {
		return tzc.data[offset : offset+tzc.arr[id]]
	}
	if id <= uint(len(tzc.arr)) {
		return tzc.data[offset+tzc.arr[id-1] : offset+tzc.arr[id]]
	}
	return nil
}

func (tzc *Timezonecache) Search(lat, lng float64) (Result, error) {
	var name string
	start := time.Now()
	ll := geo.NewLatLng(lat, lng)
	if !ll.Valid() {
		return Result{}, ErrCoordinatesNotValid
	}

	tzc.rt.SearchLatLng(ll, func(min geo.LatLng, max geo.LatLng, value interface{}) bool {
		if id, ok := value.(uint); ok {
			p := geo.NewPolygon()
			p.FromByteSlice(tzc.buf(id))
			if p.ContainsLatLng(ll) {
				name = tzc.name[id]
				return true
			}
		}
		return false
	})
	return Result{Name: name, Coordinates: ll, Elapsed: time.Since(start)}, nil
}

// Result is a timezone lookup result
type Result struct {
	Name        string
	Coordinates geo.LatLng
	Elapsed     time.Duration
}

func (tzc *Timezonecache) Save(filename string) error {
	f2, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.FileMode(0666))
	if err != nil {
		return err
	}
	defer f2.Close()

	bw := bufio.NewWriter(f2)
	buf := make([]byte, 256)
	var n, written int

	if n, err = bw.Write(tzc.encodeHeader(buf)); err != nil {
		return err
	}
	written += n
	for i := 0; i < len(tzc.name); i++ {
		if n, err = bw.Write(tzc.encodeItem(buf, i)); err != nil {
			return err
		}
		written += n
	}
	if err = bw.Flush(); err != nil {
		return err
	}
	offset := int64(written + pageSize - written%pageSize)

	if n, err = f2.WriteAt(tzc.data, offset); err != nil {
		return err
	}
	written += n
	return bw.Flush()
}

func (tzc *Timezonecache) encodeItem(buf []byte, i int) []byte {
	name := tzc.name[i]
	if len(buf) >= 5+len(name) {
		endian.PutUint32(buf, tzc.arr[i])
		buf[4] = uint8(len(name))
		copy(buf[5:5+len(name)], name)
	}
	return buf[:5+len(name)]
}

func (tzc *Timezonecache) encodeHeader(b []byte) []byte {
	headerLength := 10
	if len(b) >= 10 {
		for i := 0; i < len(tzc.name); i++ {
			headerLength += 5 + len(tzc.name)
		}
		endian.PutUint32(b[:4], uint32(headerLength))
		endian.PutUint32(b[4:8], uint32(len(tzc.data)))
		endian.PutUint16(b[8:10], uint16(len(tzc.arr)))
	}

	return b[:10]
}

func (tzc *Timezonecache) decodeHeader(b []byte) (n int, err error) {
	if len(b) > 10 {
		tzc.dataOffset = endian.Uint32(b[:4])
		tzc.dataLength = endian.Uint32(b[4:8])
		items := endian.Uint16(b[8:10])
		tzc.arr = make([]uint32, 0, items)
		tzc.name = make([]string, 0, items)
		return 10, nil
	}

	return 0, errors.New("error []byte insufficient for header")
}

func (tzc *Timezonecache) decodeItem(b []byte) (n int, err error) {
	if len(b) > 4 && len(b) >= (5+int(b[4])) {
		tzc.arr = append(tzc.arr, endian.Uint32(b[:4]))
		tzc.name = append(tzc.name, string(b[5:5+b[4]]))
		return (5 + int(b[4])), nil
	}
	return 0, errors.New("error []byte insufficient for an item")
}

func (tzc *Timezonecache) Load(f *os.File) (err error) {
	var d, discarded int
	var b []byte
	br := bufio.NewReader(f)
	if b, err = br.Peek(256); err != nil {
		return err
	}
	if d, err = tzc.decodeHeader(b); err != nil {
		return err
	}
	if d, err = br.Discard(d); err != nil {
		return err
	}
	discarded += d
	for i := 0; i < cap(tzc.name); i++ {
		if b, err = br.Peek(256); err != nil {
			return err
		}
		if d, err = tzc.decodeItem(b); err != nil {
			return err
		}
		if d, err = br.Discard(d); err != nil {
			return err
		}
		discarded += d
	}
	offset := int64(discarded + pageSize - discarded%pageSize)
	if tzc.data, err = mmap(f, offset, int64(tzc.dataLength)); err != nil {
		return err
	}
	tzc.BuildRtree()
	return err
}

func (tzc *Timezonecache) Close() error {
	if tzc.data != nil {
		err := munmap(tzc.data)
		tzc.data = nil
		return err
	}
	return errors.New("error timezone data is nil")
}

func (tzc *Timezonecache) BuildRtree() {
	for i, _ := range tzc.arr {
		id := uint(i)
		p := geo.NewPolygonFromBytes(tzc.buf(id))
		tzc.rt.InsertPolygon(p, id)
	}
}

func mmap(f *os.File, offset, length int64) ([]byte, error) {
	if f == nil {
		return nil, errors.New("error file not open")
	}
	return syscall.Mmap(int(f.Fd()), offset, int(length), syscall.PROT_READ, unix.MAP_SHARED)
}

func munmap(data []byte) (err error) {
	if data != nil {
		err = syscall.Munmap(data)
		data = nil
		return
	}
	return errors.New("error munmap data is nil")
}

// [8]items [8*items]offset [1*items]offsetstring

// [4]dataoffset [2]items [4]offset [1]stringlength [...]string ... []
