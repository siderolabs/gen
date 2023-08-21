// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xtesting provides a T interface wrapper around *testing.T
package xtesting

// T is an interface wrapper around *testing.T.
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
}
