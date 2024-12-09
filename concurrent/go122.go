// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.22 && !go1.24 && !nomaptypehash

package concurrent

import (
	"unsafe"
)

// _MapType is runtime.maptype from runtime/type.go.
//
//nolint:govet
type _MapType struct {
	_Type
	Key    *_Type
	Elem   *_Type
	Bucket *_Type // internal type representing a hash bucket
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher     func(unsafe.Pointer, uintptr) uintptr
	KeySize    uint8  // size of key slot
	ValueSize  uint8  // size of elem slot
	BucketSize uint16 // size of bucket
	Flags      uint32
}

// _Type is runtime._type from runtime/type.go.
//
//nolint:govet,revive
type _Type struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       uint8   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, GCData is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	GCData    *byte
	Str       int32 // string form
	PtrToThis int32 // type for pointer to this type, may be zero
}

// efaceMap is runtime.eface from runtime/runtime2.go.
type efaceMap struct {
	_type *_MapType
	data  unsafe.Pointer
}

func efaceMapOf(ep any) *efaceMap {
	return (*efaceMap)(unsafe.Pointer(&ep))
}
