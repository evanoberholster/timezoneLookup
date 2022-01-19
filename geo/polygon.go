// Copyright 2022 Evan Oberholster. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geo

import (
	"reflect"
	"unsafe"
)

// Polygon represents a closed Polygon of vertices when
// the first and last vertices are equal.
type Polygon struct {
	min, max LatLng   // min and man LatLng
	v        []LatLng // Vertices
}

// NewPolygon returns a new empty Polygon
func NewPolygon() Polygon {
	return Polygon{
		max: NewLatLng(minLatitude, minLongitude),
		min: NewLatLng(maxLatitude, maxLongitude),
	}
}

// NewPolygonFromVertices returns a new Polygon with the LatLng vertices.
// updates the polygon's boundingbox.
func NewPolygonFromVertices(s []LatLng) Polygon {
	p := NewPolygon()
	p.v = s
	p.UpdateBoundingBox()
	return p
}

func NewPolygonFromBytes(b []byte) Polygon {
	p := NewPolygon()
	p.v = toLatLngSlice(b)
	p.UpdateBoundingBox()
	return p
}

// updateBounds updates the max and min limits of the boundingBox.
func (p *Polygon) updateBounds(ll LatLng) {
	if ll.Valid() {
		if p.max.Lat < ll.Lat {
			p.max.Lat = ll.Lat
		}
		if p.max.Lng < ll.Lng {
			p.max.Lng = ll.Lng
		}
		if p.min.Lat > ll.Lat {
			p.min.Lat = ll.Lat
		}
		if p.min.Lng > ll.Lng {
			p.min.Lng = ll.Lng
		}
	}
}

// Length retuns the number of veritices in the Polygon
func (p *Polygon) Length() int {
	return len(p.v)
}

// Max returns the bottom-left coordinate of the Polygon.
// Correspoinding to the minimum latitide and longitude values contained.
func (p *Polygon) Min() LatLng {
	return p.min
}

// Max returns the top-right coordinate of the Polygon.
// Correspoinding to the maximum latitude and longitude values contained.
func (p *Polygon) Max() LatLng {
	return p.max
}

// Add adds a Latitude and Longitude in degrees to the Polygon.
// Maximum and minimum latitude is -90 and +90 respectively.
// Maximum and minimum longitude is -180 and +180 respectively.
func (p *Polygon) Add(latitude, longitude float64) {
	p.AddVertex(LatLng{float32(latitude), float32(longitude)})
}

// AddVertex adds a LatLng Vertex to the Polygon. Updates polygon bounds with new Vertex.
func (p *Polygon) AddVertex(ll LatLng) {
	if ll.Valid() {
		p.v = append(p.v, ll)
		p.updateBounds(ll)
	}
}

// UpdateBoundingBox updates the max and min limits of the boundingBox using the contained ploygon vertices.
func (p *Polygon) UpdateBoundingBox() {
	for _, v := range p.v {
		p.updateBounds(v)
	}
}

func (p *Polygon) ContainsLatLng(query LatLng) bool {
	if len(p.v) < 3 {
		return false
	}
	in := rayIntersectsSegment(query, p.v[len(p.v)-1], p.v[0])
	for i := 1; i < len(p.v); i++ {
		if rayIntersectsSegment(query, p.v[i-1], p.v[i]) {
			in = !in
		}
	}
	return in
}

func rayIntersectsSegment(p, a, b LatLng) bool {
	return (a.Lng > p.Lng) != (b.Lng > p.Lng) &&
		p.Lat < (b.Lat-a.Lat)*(p.Lng-a.Lng)/(b.Lng-a.Lng)+a.Lat
}

// reference: https://go101.org/article/unsafe.html
func toByteSlice(b []LatLng) []byte {
	var bs []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	hdr.Len = len(b) * 8
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	return bs
}

// reference: https://go101.org/article/unsafe.html
func toLatLngSlice(b []byte) (result []LatLng) {
	var lls []LatLng
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&lls))
	hdr.Len = len(b) / 8
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	return lls
}

func (p Polygon) ToByteSlice() []byte {
	return toByteSlice(p.v)
}

func (p *Polygon) FromByteSlice(src []byte) {
	p.v = toLatLngSlice(src)
}
