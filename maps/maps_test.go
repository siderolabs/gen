// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package maps_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/siderolabs/gen/maps"
)

func TestFilterInPlace(t *testing.T) {
	t.Parallel()

	type args struct {
		m map[string]string
	}

	tests := map[string]struct {
		args args
		want map[string]string
	}{
		"nil": {
			args: args{
				m: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				m: map[string]string{},
			},
			want: map[string]string{},
		},
		"single": {
			args: args{
				m: map[string]string{"foo": "b"},
			},
			want: map[string]string{"foo": "b"},
		},
		"multiple": {
			args: args{
				m: map[string]string{"foo": "b", "bar": "c", "baz": "d"},
			},
			want: map[string]string{"foo": "b"},
		},
		"multiple to empty": {
			args: args{
				m: map[string]string{"far": "b", "bar": "c", "baz": "d"},
			},
			want: map[string]string{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.FilterInPlace(tt.args.m, func(k, _ string) bool { return k == "foo" })

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	type args struct {
		m map[string]string
	}

	tests := map[string]struct {
		args args
		want map[string]string
	}{
		"nil": {
			args: args{
				m: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				m: map[string]string{},
			},
			want: nil,
		},
		"single": {
			args: args{
				m: map[string]string{"foo": "b"},
			},
			want: map[string]string{"foo": "b"},
		},
		"multiple": {
			args: args{
				m: map[string]string{"foo": "b", "bar": "c", "baz": "d"},
			},
			want: map[string]string{"foo": "b"},
		},
		"multiple to empty": {
			args: args{
				m: map[string]string{"far": "b", "bar": "c", "baz": "d"},
			},
			want: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.Filter(tt.args.m, func(k, _ string) bool { return k == "foo" })

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()

	type args struct {
		m map[string]string
	}

	tests := map[string]struct {
		args args
		want []string
	}{
		"nil": {
			args: args{
				m: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				m: map[string]string{},
			},
			want: nil,
		},
		"single": {
			args: args{
				m: map[string]string{"foo": "b"},
			},
			want: []string{"foo"},
		},
		"multiple": {
			args: args{
				m: map[string]string{"foo": "b", "bar": "c", "baz": "d"},
			},
			want: []string{"bar", "baz", "foo"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.Keys(tt.args.m)

			slices.Sort(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKeysFunc(t *testing.T) {
	t.Parallel()

	type args struct {
		m map[string]string
	}

	tests := map[string]struct {
		args args
		want []string
	}{
		"nil": {
			args: args{
				m: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				m: map[string]string{},
			},
			want: nil,
		},
		"single": {
			args: args{
				m: map[string]string{"foo": "b"},
			},
			want: []string{"foo func"},
		},
		"multiple": {
			args: args{
				m: map[string]string{"foo": "b", "bar": "c", "baz": "d"},
			},
			want: []string{"bar func", "baz func", "foo func"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.KeysFunc(tt.args.m, func(k string) string { return k + " func" })

			slices.Sort(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToSlice(t *testing.T) {
	t.Parallel()

	type args struct {
		m map[string]string
	}

	tests := map[string]struct {
		args args
		want []string
	}{
		"nil": {
			args: args{
				m: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				m: map[string]string{},
			},
			want: nil,
		},
		"single": {
			args: args{
				m: map[string]string{"foo": "b"},
			},
			want: []string{"foo b"},
		},
		"multiple": {
			args: args{
				m: map[string]string{"foo": "b", "bar": "c", "baz": "d"},
			},
			want: []string{"bar c", "baz d", "foo b"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.ToSlice(tt.args.m, func(k, v string) string { return k + " " + v })

			slices.Sort(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValuesFunc(t *testing.T) {
	t.Parallel()

	type args struct {
		m map[string]string
	}

	tests := map[string]struct {
		args args
		want []string
	}{
		"nil": {
			args: args{
				m: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				m: map[string]string{},
			},
			want: nil,
		},
		"single": {
			args: args{
				m: map[string]string{"foo": "b"},
			},
			want: []string{"b"},
		},
		"multiple": {
			args: args{
				m: map[string]string{"foo": "b", "bar": "c", "baz": "d"},
			},
			want: []string{"b", "c", "d"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.ValuesFunc(tt.args.m, func(v string) string { return v })

			slices.Sort(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	type args struct {
		maps []map[int]struct{}
	}

	tests := map[string]struct {
		args args
		want []int
	}{
		"nil": {
			args: args{
				maps: nil,
			},
			want: nil,
		},
		"empty": {
			args: args{
				maps: []map[int]struct{}{
					{},
				},
			},
			want: nil,
		},
		"single": {
			args: args{
				maps: []map[int]struct{}{
					{
						1: {},
						2: {},
					},
				},
			},
			want: []int{1, 2},
		},
		"first empty": {
			args: args{
				maps: []map[int]struct{}{
					{},
					{
						1: {},
						2: {},
					},
				},
			},
			want: nil,
		},
		"multiple": {
			args: args{
				maps: []map[int]struct{}{
					{
						1: {},
						2: {},
						4: {},
					},
					{
						3: {},
						2: {},
						4: {},
					},
					{
						2: {},
						4: {},
					},
				},
			},
			want: []int{2, 4},
		},
		"empty intersection": {
			args: args{
				maps: []map[int]struct{}{
					{
						4: {},
						5: {},
						6: {},
					},
					{
						5: {},
						6: {},
					},
					{
						1: {},
						2: {},
						3: {},
					},
				},
			},
			want: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := maps.Intersect(tt.args.maps...)

			slices.Sort(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKeysAddtional(t *testing.T) {
	m := generateMap(6)

	keys := maps.Keys(m)

	assert.Equal(t, 6, len(keys))
	slices.Sort(keys)
	assert.EqualValues(t, []int{0, 1, 2, 3, 4, 5}, keys)
}

func TestValuesAddtional(t *testing.T) {
	m := generateMap(6)

	values := maps.Values(m)

	assert.Equal(t, 6, len(values))
	slices.Sort(values)
	assert.EqualValues(t, []int{-5, -4, -3, -2, -1, 0}, values)
}

var Sink []int

func BenchmarkKeys(b *testing.B) {
	smallMap := generateMap(10)
	midMap := generateMap(100)
	largeMap := generateMap(1000)

	b.Run("small", func(b *testing.B) {
		for range b.N {
			Sink = maps.Keys(smallMap)
		}
	})

	b.Run("mid", func(b *testing.B) {
		for range b.N {
			Sink = maps.Keys(midMap)
		}
	})

	b.Run("large", func(b *testing.B) {
		for range b.N {
			Sink = maps.Keys(largeMap)
		}
	})
}

func BenchmarkValues(b *testing.B) {
	smallMap := generateMap(10)
	midMap := generateMap(100)
	largeMap := generateMap(1000)

	b.Run("small", func(b *testing.B) {
		for range b.N {
			Sink = maps.Values(smallMap)
		}
	})

	b.Run("mid", func(b *testing.B) {
		for range b.N {
			Sink = maps.Values(midMap)
		}
	})

	b.Run("large", func(b *testing.B) {
		for range b.N {
			Sink = maps.Values(largeMap)
		}
	})
}

func generateMap(num int) map[int]int {
	result := make(map[int]int, num)

	for i := range num {
		result[i] = -i
	}

	return result
}
