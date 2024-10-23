// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xstrings_test

import (
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/siderolabs/gen/xiter/xstrings"
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
		result := slices.Collect(xstrings.Lines(s.a))
		if !slices.Equal(result, s.b) {
			t.Errorf(`slices.Collect(Lines(%q)) = %q; want %q`, s.a, result, s.b)
		}
	}
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
	{abcd, "", 2, []string{"a", "bcd"}},
	{abcd, "", 4, []string{"a", "b", "c", "d"}},
	{abcd, "", -1, []string{"a", "b", "c", "d"}},
	{faces, "", -1, []string{"☺", "☻", "☹"}},
	{faces, "", 3, []string{"☺", "☻", "☹"}},
	{faces, "", 17, []string{"☺", "☻", "☹"}},
	{"☺�☹", "", -1, []string{"☺", "�", "☹"}},
	{abcd, "a", 0, nil},
	{abcd, "a", -1, []string{"", "bcd"}},
	{abcd, "z", -1, []string{"abcd"}},
	{commas, ",", -1, []string{"1", "2", "3", "4"}},
	{dots, "...", -1, []string{"1", ".2", ".3", ".4"}},
	{faces, "☹", -1, []string{"☺☻", ""}},
	{faces, "~", -1, []string{faces}},
	{"1 2 3 4", " ", 3, []string{"1", "2", "3 4"}},
	{"1 2", " ", 3, []string{"1", "2"}},
	{"", "T", math.MaxInt / 4, []string{""}},
	{"\xff-\xff", "", -1, []string{"\xff", "-", "\xff"}},
	{"\xff-\xff", "-", -1, []string{"\xff", "\xff"}},
}

func TestSplit(t *testing.T) {
	for _, tt := range splittests {
		a := strings.SplitN(tt.s, tt.sep, tt.n)
		if !slices.Equal(a, tt.a) {
			t.Errorf("Split(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, a, tt.a)

			continue
		}

		if tt.n < 0 {
			a2 := slices.Collect(xstrings.SplitSeq(tt.s, tt.sep))
			if !slices.Equal(a2, tt.a) {
				t.Errorf(`collect(SplitSeq(%q, %q)) = %v; want %v`, tt.s, tt.sep, a2, tt.a)
			}
		}

		if tt.n == 0 {
			continue
		}

		s := strings.Join(a, tt.sep)
		if s != tt.s {
			t.Errorf("Join(Split(%q, %q, %d), %q) = %q", tt.s, tt.sep, tt.n, tt.sep, s)
		}

		if tt.n < 0 {
			b := strings.Split(tt.s, tt.sep)
			if !slices.Equal(a, b) {
				t.Errorf("Split disagrees with SplitN(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, b, a)
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
		a := strings.SplitAfterN(tt.s, tt.sep, tt.n)
		if !slices.Equal(a, tt.a) {
			t.Errorf(`Split(%q, %q, %d) = %v; want %v`, tt.s, tt.sep, tt.n, a, tt.a)

			continue
		}

		if tt.n < 0 {
			a2 := slices.Collect(xstrings.SplitAfterSeq(tt.s, tt.sep))
			if !slices.Equal(a2, tt.a) {
				t.Errorf(`collect(SplitAfterSeq(%q, %q)) = %v; want %v`, tt.s, tt.sep, a2, tt.a)
			}
		}

		s := strings.Join(a, "")
		if s != tt.s {
			t.Errorf(`Join(Split(%q, %q, %d), %q) = %q`, tt.s, tt.sep, tt.n, tt.sep, s)
		}

		if tt.n < 0 {
			b := strings.SplitAfter(tt.s, tt.sep)
			if !slices.Equal(a, b) {
				t.Errorf("SplitAfter disagrees with SplitAfterN(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, b, a)
			}
		}
	}
}
