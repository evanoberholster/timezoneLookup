// Copyright 2022 Evan Oberholster. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geo

const (
	minLatitude  = -90
	maxLatitude  = 90
	minLongitude = -180
	maxLongitude = 180
)

// LatLng represents Latitude and Longitude in degree format
type LatLng struct {
	Lat, Lng float32
}

func (ll LatLng) toFloat32() [2]float32 {
	return [2]float32{ll.Lat, ll.Lng}
}

// Valid returns true of LatLng is with the bounds of minLatitude/maxLatitude and minLongitude/maxLongitude
func (ll LatLng) Valid() bool {
	return (ll.Lat >= minLatitude && ll.Lat <= maxLatitude) && (ll.Lng >= minLongitude && ll.Lng <= maxLongitude)
}

// NewLatLng returns a new LatLng with the given latitude and longitude
func NewLatLng(latitude, longitude float64) LatLng {
	return LatLng{float32(latitude), float32(longitude)}
}

// SearchLatLng searchs the RTree for the given LatLng combination
func (tr *RTree) SearchLatLng(ll LatLng, iter func(min, max LatLng, value interface{}) bool) {
	tr.searchLatLng(rect{min: ll.toFloat32(), max: ll.toFloat32()}, iter)
}

// InsertPolygon data into tree
func (tr *RTree) InsertPolygon(p Polygon, value interface{}) {
	var item rect
	fit(p.min.toFloat32(), p.max.toFloat32(), value, &item)
	tr.insert(&item)
}

func (tr *RTree) searchLatLng(
	target rect,
	iter func(min, max LatLng, value interface{}) bool,
) {
	if tr.root.data == nil {
		return
	}
	if target.intersects(&tr.root) {
		tr.root.searchLatLng(target, tr.height, iter)
	}
}

func (r *rect) searchLatLng(
	target rect, height int,
	iter func(min, max LatLng, value interface{}) bool,
) bool {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if target.intersects(&n.rects[i]) {
				if !iter(LatLng{n.rects[i].min[0], n.rects[i].min[1]}, LatLng{n.rects[i].max[0], n.rects[i].max[1]}, n.rects[i].data) {
					return false
				}
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if target.intersects(&n.rects[i]) {
				if !n.rects[i].searchLatLng(target, height-1, iter) {
					return false
				}
			}
		}
	}
	return true
}
