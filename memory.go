// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package timezoneLookup

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/golang/snappy"
)

type Memory struct { // Memory struct
	filename  string
	timezones []Timezone
	snappy    bool
}

func MemoryStorage(snappy bool, filename string) *Memory {
	if snappy {
		filename = filename + ".snap.json"
	} else {
		filename = filename + ".json"
	}
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
	data, err := json.Marshal(m.timezones)
	if err != nil {
		return err
	}
	w, err := os.Create(dbFilename)
	if err != nil {
		return err
	}
	defer w.Close()
	if m.snappy {
		snap := snappy.NewBufferedWriter(w)
		_, err := snap.Write(data)
		if err != nil {
			return err
		}
		defer snap.Close()
	} else {
		_, err = w.Write(data)
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
