// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xstrings provides additional iterator functions for working with strings.
// Functions from this package will become a part of the Go standard library in Go 1.24, when this package
// functions will be deprecated and eventually removed.
package xstrings

import (
	"iter"
	"strings"
	"unicode/utf8"
)

// Lines returns an iterator over the newline-terminated lines in the string s.
// The lines yielded by the iterator include their terminating newlines.
// If s is empty, the iterator yields no lines at all.
// If s does not end in a newline, the final yielded line will not end in a newline.
// It returns a single-use iterator.
func Lines(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(s) > 0 {
			var line string

			if i := strings.IndexByte(s, '\n'); i >= 0 {
				line, s = s[:i+1], s[i+1:]
			} else {
				line, s = s, ""
			}

			if !yield(line) {
				return
			}
		}
	}
}

// explodeSeq returns an iterator over the runes in s.
func explodeSeq(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(s) > 0 {
			_, size := utf8.DecodeRuneInString(s)
			if !yield(s[:size]) {
				return
			}

			s = s[size:]
		}
	}
}

// splitSeq is SplitSeq or SplitAfterSeq, configured by how many
// bytes of sep to include in the results (none or all).
func splitSeq(s, sep string, sepSave int) iter.Seq[string] {
	if len(sep) == 0 {
		return explodeSeq(s)
	}

	return func(yield func(string) bool) {
		for {
			i := strings.Index(s, sep)
			if i < 0 {
				break
			}

			frag := s[:i+sepSave]
			if !yield(frag) {
				return
			}

			s = s[i+len(sep):]
		}

		yield(s)
	}
}

// SplitSeq returns an iterator over all substrings of s separated by sep.
// The iterator yields the same strings that would be returned by Split(s, sep),
// but without constructing the slice.
// It returns a single-use iterator.
func SplitSeq(s, sep string) iter.Seq[string] {
	return splitSeq(s, sep, 0)
}

// SplitAfterSeq returns an iterator over substrings of s split after each instance of sep.
// The iterator yields the same strings that would be returned by SplitAfter(s, sep),
// but without constructing the slice.
// It returns a single-use iterator.
func SplitAfterSeq(s, sep string) iter.Seq[string] {
	return splitSeq(s, sep, len(sep))
}
