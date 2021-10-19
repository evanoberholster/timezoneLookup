// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package timezoneLookup

import (
	"errors"
	"io"
	"os"

	json "github.com/goccy/go-json"
	"github.com/klauspost/compress/snappy"
)

type Memory struct { // Memory struct
	filename  string
	timezones []Timezone
	snappy    bool
}

func MemoryStorage(snappy bool, filename string) *Memory {
	if snappy {
		filename += ".snap"
	}
	filename += ".json"
	return &Memory{
		filename:  filename,
		timezones: []Timezone{},
		snappy:    snappy,
	}
}

func (m *Memory) Close() {
	m.timezones = []Timezone{}
}

func (m *Memory) LoadTimezones() error {
	file, err := os.Open(m.filename)
	if err != nil {
		return err
	}

	var tzs []Timezone
	if m.snappy {
		data := snappy.NewReader(file)
		dec := json.NewDecoder(data)
		for dec.More() {
			err := dec.Decode(&tzs)
			if err != nil {
				return err
			}
		}
	} else {
		dec := json.NewDecoder(file)
		for dec.More() {

			err := dec.Decode(&tzs)
			if err != nil {
				return err
			}
		}
	}

	m.timezones = tzs
	return nil
}

func (m *Memory) Query(q Coord) (string, error) {
	for _, tz := range m.timezones {
		for _, p := range tz.Polygons {
			if p.Min.Lat < q.Lat && p.Min.Lon < q.Lon && p.Max.Lat > q.Lat && p.Max.Lon > q.Lon {
				if p.contains(q) {
					return tz.Tzid, nil
				}
			}
		}
	}
	return "Error", errors.New(errTimezoneNotFound)
}

func (m *Memory) writeTimezoneJSON(dbFilename string) error {
	w, err := os.Create(dbFilename)
	if err != nil {
		return err
	}
	defer w.Close()
	sw := io.WriteCloser(w)
	if m.snappy {
		sw = snappy.NewBufferedWriter(w)
	}
	err = json.NewEncoder(sw).Encode(m.timezones)
	if closeErr := sw.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if m.snappy {
		if closeErr := w.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}

func (m *Memory) CreateTimezones(jsonFilename string) error {
	tzs, err := TimezonesFromGeoJSON(jsonFilename)
	if err != nil {
		return err
	}
	m.timezones = tzs
	err = m.writeTimezoneJSON(m.filename)
	if err != nil {
		return err
	}
	return nil
}
