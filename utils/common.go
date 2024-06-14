// Copyright 2021-2024 IBM Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var errInvalidByteQuantity = errors.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")
var errInvalidBoolQuantity = errors.New("bool quantity must be a string like ON, OFF, YES, NO")

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

// ToBytes parses a string formatted by ByteSize as bytes. Note binary-prefixed and SI prefixed units both mean a base-2 units
// KB = K = KiB	= 1024
// MB = M = MiB = 1024 * K
// GB = G = GiB = 1024 * M
// TB = T = TiB = 1024 * G
func ToBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return 0, errInvalidByteQuantity
	}

	bytesString, multiple := s[:i], s[i:]
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes < 0 {
		return 0, errInvalidByteQuantity
	}

	switch multiple {
	case "T", "TB", "TIB":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GIB":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MIB":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KIB":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, errInvalidByteQuantity
	}
}

// ToBool parses a string formatted to bool, but prometheus only accepts float64 formatted data, so this function retuns float64 formatted.
func ToBool(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return -1, errInvalidBoolQuantity
	}

	switch s {
	case "ON", "YES":
		return 1, nil
	case "OFF", "NO":
		return 0, nil
	default:
		return -1, errInvalidBoolQuantity
	}
}
