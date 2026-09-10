// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xyaml_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"

	"github.com/siderolabs/gen/xyaml"
)

type leaf struct {
	Name string `yaml:"name"`
}

type nested struct {
	Ptrs   []*leaf          `yaml:"ptrs"`
	Values []leaf           `yaml:"values"`
	Map    map[string]*leaf `yaml:"map"`
	Any    []any            `yaml:"any"`
	Inner  *nested          `yaml:"inner"`
	Args   argValue         `yaml:"args"`
}

type rootDoc struct {
	Machine *nested `yaml:"machine"`
}

type inlineDoc struct {
	Extra map[string]*leaf `yaml:",inline"`
}

func TestCheckNullElements(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		data string
		errs []string
	}{
		{
			name: "valid",
			data: `
machine:
  ptrs:
    - name: a
    - name: b
  values:
    - name: c
  map:
    x:
      name: d
  inner:
    ptrs:
      - name: e
`,
		},
		{
			name: "null struct field is allowed",
			data: `
machine:
  inner: null
`,
		},
		{
			name: "null value element is allowed",
			data: `
machine:
  values:
    - null
    - ~
`,
		},
		{
			name: "null any element is allowed",
			data: `
machine:
  any:
    - null
    - ~
`,
		},
		{
			name: "null pointer element",
			data: `
machine:
  ptrs:
    - name: a
    - null
`,
			errs: []string{
				`null value is not allowed at "machine.ptrs[1]" (line 5)`,
			},
		},
		{
			name: "empty pointer element",
			data: `
machine:
  ptrs:
    -
    - ~
`,
			errs: []string{
				`null value is not allowed at "machine.ptrs[0]" (line 4)`,
				`null value is not allowed at "machine.ptrs[1]" (line 5)`,
			},
		},
		{
			name: "null map value",
			data: `
machine:
  map:
    x: null
    y:
      name: b
`,
			errs: []string{
				`null value is not allowed at "machine.map.x" (line 4)`,
			},
		},
		{
			name: "nested",
			data: `
machine:
  inner:
    ptrs:
        - null
`,
			errs: []string{
				`null value is not allowed at "machine.inner.ptrs[0]" (line 5)`,
			},
		},
		{
			name: "alias to null",
			data: `
machine:
  inner: &n null
  ptrs:
    - *n
`,
			errs: []string{
				`null value is not allowed at "machine.ptrs[0]" (line 3)`,
			},
		},
		{
			name: "alias to sequence with null",
			data: `
machine:
  any: &l
    - null
  ptrs: *l
`,
			errs: []string{
				`null value is not allowed at "machine.ptrs[0]" (line 4)`,
			},
		},
		{
			name: "alias to mapping with null",
			data: `
machine:
  any:
    - &m
      ptrs:
        - null
  inner: *m
`,
			errs: []string{
				`null value is not allowed at "machine.inner.ptrs[0]" (line 6)`,
			},
		},
		{
			name: "alias to valid mapping",
			data: `
machine:
  inner: &m
    ptrs:
      - name: a
  ptrs:
    - *m
`,
		},
		{
			name: "unknown keys are skipped",
			data: `
machine:
  unknown:
    - null
`,
		},
		{
			name: "custom unmarshaler",
			data: `
machine:
  args:
    - null
`,
		},
		{
			name: "empty document",
			data: ``,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var node yaml.Node

			require.NoError(t, yaml.Unmarshal([]byte(tt.data), &node))

			err := xyaml.CheckNullElements(&rootDoc{}, &node)

			if len(tt.errs) == 0 {
				require.NoError(t, err)

				return
			}

			require.Error(t, err)

			for _, expected := range tt.errs {
				require.ErrorContains(t, err, expected)
			}

			require.Len(t, splitErrors(err), len(tt.errs))
		})
	}
}

func TestCheckNullElementsInlineMap(t *testing.T) {
	t.Parallel()

	var node yaml.Node

	require.NoError(t, yaml.Unmarshal([]byte("a:\n  name: x\nb: null\n"), &node))

	err := xyaml.CheckNullElements(&inlineDoc{}, &node)
	require.ErrorContains(t, err, `null value is not allowed at "b" (line 3)`)
	require.Len(t, splitErrors(err), 1)
}

func TestCheckNullElementsAliasCycle(t *testing.T) {
	var node yaml.Node

	require.NoError(t, yaml.Unmarshal([]byte("machine: &m\n  inner: *m\n"), &node))

	require.ErrorContains(t, xyaml.CheckNullElements(&rootDoc{}, &node), `anchor "m" value contains itself`)
}

func TestCheckNullElementsTypeMismatch(t *testing.T) {
	t.Parallel()

	var node yaml.Node

	require.NoError(t, yaml.Unmarshal([]byte("machine:\n  values: {a: b}\n"), &node))

	require.ErrorContains(t, xyaml.CheckNullElements(&rootDoc{}, &node), "unexpected type for yaml mapping")
}

func splitErrors(err error) []error {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		return joined.Unwrap()
	}

	return []error{err}
}
