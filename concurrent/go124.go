// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.24 && !go1.25 && !nomaptypehash

package concurrent

import (
	"unsafe"
)

// _MapType is runtime.maptype [abi.SwissMapType] from runtime/type.go.
//
//nolint:govet
type _MapType struct {
	_Type
	Key   *_Type
	Elem  *_Type
	Group *_Type // internal type representing a slot group
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher    func(unsafe.Pointer, uintptr) uintptr
	GroupSize uintptr // == Group.Size_
	SlotSize  uintptr // size of key/elem slot
	ElemOff   uintptr // offset of elem in key/elem slot
	Flags     uint32
}

// _Type is runtime._type from runtime/type.go.
//
//nolint:govet,revive
type _Type struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       _TFlag  // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       _Kind   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// Normally, GCData points to a bitmask that describes the
	// ptr/nonptr fields of the type. The bitmask will have at
	// least PtrBytes/ptrSize bits.
	// If the TFlagGCMaskOnDemand bit is set, GCData is instead a
	// **byte and the pointer to the bitmask is one dereference away.
	// The runtime will build the bitmask if needed.
	// (See runtime/type.go:getGCMask.)
	// Note: multiple types may have the same value of GCData,
	// including when TFlagGCMaskOnDemand is set. The types will, of course,
	// have the same pointer layout (but not necessarily the same size).
	GCData    *byte
	Str       _NameOff // string form
	PtrToThis _TypeOff // type for pointer to this type, may be zero
}

// _TypeOff is the offset to a type from moduledata.types.  See resolveTypeOff in runtime.
type _TypeOff int32

// _NameOff is the offset to a name from moduledata.types.  See resolveNameOff in runtime.
type _NameOff int32

type _Kind uint8

type _TFlag uint8

// efaceMap is runtime.eface from runtime/runtime2.go.
type efaceMap struct {
	_type *_MapType
	data  unsafe.Pointer
}

func efaceMapOf(ep any) *efaceMap {
	return (*efaceMap)(unsafe.Pointer(&ep))
}
