# Copyright 2021 Tamás Gulácsi. 
#
# SPDX-License-Identifier: MIT

using Go = import "/go.capnp";
$Go.package("cp");
$Go.import("github.com/evanoberholster/timezoneLookup/cp");
@0xcc3309152ef8d9f0;

struct Coord @0xa6c71d8345916b7b {
    lat @0 :Float32;
    lon @1 :Float32;
}

struct Polygon @0xb1b308e89f4237bf {
    max @0 :Coord;
    min @1 :Coord;
    coords @2 :List(Coord);
}

struct PolygonIndex @0xf7e516726966c831 {
    id @0 :UInt64;
    max @1 :Coord;
    min @2 :Coord;
    tzid @3 :Text;
}

