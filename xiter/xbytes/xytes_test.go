// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xbytes_test

import (
	"bytes"
	"math"
	"slices"
	"testing"

	"github.com/siderolabs/gen/xiter/xbytes"
)

type LinesTest struct {
	a string
	b []string
}

//nolint:dupword
var linesTests = []LinesTest{
	{a: "abc\nabc\n", b: []string{"abc\n", "abc\n"}},
	{a: "abc\r\nabc", b: []string{"abc\r\n", "abc"}},
	{a: "abc\r\n", b: []string{"abc\r\n"}},
	{a: "\nabc", b: []string{"\n", "abc"}},
	{a: "\nabc\n\n", b: []string{"\n", "abc\n", "\n"}},
}

func TestLines(t *testing.T) {
	for _, s := range linesTests {
		result := sliceOfString(slices.Collect(xbytes.Lines([]byte(s.a))))
		if !slices.Equal(result, s.b) {
			t.Errorf(`slices.Collect(Lines(%q)) = %q; want %q`, s.a, result, s.b)
		}
	}
}

func sliceOfString(s [][]byte) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = string(v)
	}

	return result
}

//nolint:govet
type SplitTest struct {
	s   string
	sep string
	n   int
	a   []string
}

const (
	abcd   = "abcd"
	faces  = "☺☻☹"
	commas = "1,2,3,4"
	dots   = "1....2....3....4"
)

var splittests = []SplitTest{
	{"", "", -1, []string{}},
	{abcd, "a", 0, nil},
	{abcd, "", 2, []string{"a", "bcd"}},
	{abcd, "a", -1, []string{"", "bcd"}},
	{abcd, "z", -1, []string{"abcd"}},
	{abcd, "", -1, []string{"a", "b", "c", "d"}},
	{commas, ",", -1, []string{"1", "2", "3", "4"}},
	{dots, "...", -1, []string{"1", ".2", ".3", ".4"}},
	{faces, "☹", -1, []string{"☺☻", ""}},
	{faces, "~", -1, []string{faces}},
	{faces, "", -1, []string{"☺", "☻", "☹"}},
	{"1 2 3 4", " ", 3, []string{"1", "2", "3 4"}},
	{"1 2", " ", 3, []string{"1", "2"}},
	{"123", "", 2, []string{"1", "23"}},
	{"123", "", 17, []string{"1", "2", "3"}},
	{"bT", "T", math.MaxInt / 4, []string{"b", ""}},
	{"\xff-\xff", "", -1, []string{"\xff", "-", "\xff"}},
	{"\xff-\xff", "-", -1, []string{"\xff", "\xff"}},
}

func TestSplit(t *testing.T) {
	for _, tt := range splittests {
		a := bytes.SplitN([]byte(tt.s), []byte(tt.sep), tt.n)

		// Appending to the results should not change future results.
		var x []byte
		for _, v := range a {
			x = append(v, 'z') //nolint:gocritic
		}

		result := sliceOfString(a)
		if !slices.Equal(result, tt.a) {
			t.Errorf(`Split(%q, %q, %d) = %v; want %v`, tt.s, tt.sep, tt.n, result, tt.a)

			continue
		}

		if tt.n < 0 {
			b := sliceOfString(slices.Collect(xbytes.SplitSeq([]byte(tt.s), []byte(tt.sep))))
			if !slices.Equal(b, tt.a) {
				t.Errorf(`collect(SplitSeq(%q, %q)) = %v; want %v`, tt.s, tt.sep, b, tt.a)
			}
		}

		if tt.n == 0 || len(a) == 0 {
			continue
		}

		if want := tt.a[len(tt.a)-1] + "z"; string(x) != want {
			t.Errorf("last appended result was %s; want %s", x, want)
		}

		s := bytes.Join(a, []byte(tt.sep))
		if string(s) != tt.s {
			t.Errorf(`Join(Split(%q, %q, %d), %q) = %q`, tt.s, tt.sep, tt.n, tt.sep, s)
		}

		if tt.n < 0 {
			b := sliceOfString(bytes.Split([]byte(tt.s), []byte(tt.sep)))
			if !slices.Equal(result, b) {
				t.Errorf("Split disagrees withSplitN(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, b, a)
			}
		}

		if len(a) > 0 {
			in, out := a[0], s
			if cap(in) == cap(out) && &in[:1][0] == &out[:1][0] {
				t.Errorf("Join(%#v, %q) didn't copy", a, tt.sep)
			}
		}
	}
}

var splitaftertests = []SplitTest{
	{abcd, "a", -1, []string{"a", "bcd"}},
	{abcd, "z", -1, []string{"abcd"}},
	{abcd, "", -1, []string{"a", "b", "c", "d"}},
	{commas, ",", -1, []string{"1,", "2,", "3,", "4"}},
	{dots, "...", -1, []string{"1...", ".2...", ".3...", ".4"}},
	{faces, "☹", -1, []string{"☺☻☹", ""}},
	{faces, "~", -1, []string{faces}},
	{faces, "", -1, []string{"☺", "☻", "☹"}},
	{"1 2 3 4", " ", 3, []string{"1 ", "2 ", "3 4"}},
	{"1 2 3", " ", 3, []string{"1 ", "2 ", "3"}},
	{"1 2", " ", 3, []string{"1 ", "2"}},
	{"123", "", 2, []string{"1", "23"}},
	{"123", "", 17, []string{"1", "2", "3"}},
}

func TestSplitAfter(t *testing.T) {
	for _, tt := range splitaftertests {
		a := bytes.SplitAfterN([]byte(tt.s), []byte(tt.sep), tt.n)

		// Appending to the results should not change future results.
		var x []byte
		for _, v := range a {
			x = append(v, 'z') //nolint:gocritic
		}

		result := sliceOfString(a)
		if !slices.Equal(result, tt.a) {
			t.Errorf(`Split(%q, %q, %d) = %v; want %v`, tt.s, tt.sep, tt.n, result, tt.a)

			continue
		}

		if tt.n < 0 {
			b := sliceOfString(slices.Collect(xbytes.SplitAfterSeq([]byte(tt.s), []byte(tt.sep))))
			if !slices.Equal(b, tt.a) {
				t.Errorf(`collect(SplitAfterSeq(%q, %q)) = %v; want %v`, tt.s, tt.sep, b, tt.a)
			}
		}

		if want := tt.a[len(tt.a)-1] + "z"; string(x) != want {
			t.Errorf("last appended result was %s; want %s", x, want)
		}

		s := bytes.Join(a, nil)
		if string(s) != tt.s {
			t.Errorf(`Join(Split(%q, %q, %d), %q) = %q`, tt.s, tt.sep, tt.n, tt.sep, s)
		}

		if tt.n < 0 {
			b := sliceOfString(bytes.SplitAfter([]byte(tt.s), []byte(tt.sep)))
			if !slices.Equal(result, b) {
				t.Errorf("SplitAfter disagrees withSplitAfterN(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, b, a)
			}
		}
	}
}
