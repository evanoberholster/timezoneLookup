//go:build example
// +build example

// Copyright 2018-2022 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	timezone "github.com/evanoberholster/timezoneLookup@v2.0.0"
)

func main() {
	var tzc timezone.Timezonecache
	f, err := os.Open("timzone.data")
	if err != nil {
		return timezone.Result{}, err
	}
	defer f.Close()
	if err = tzc.Load(f); err != nil {
		return timezone.Result{}, err
	}
	defer tzc.Close()

	result := tzc.Search(lat, lng)
	fmt.Println(result)
}
