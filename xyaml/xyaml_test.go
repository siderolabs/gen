// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xyaml_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xyaml"
)

type A struct {
	Field string            `yaml:"field"`
	Map   map[string]string `yaml:"map"`
	Slice []A               `yaml:"slice"`
}

//go:embed testdata/valid.yaml
var valid []byte

//go:embed testdata/invalid.yaml
var invalid []byte

//go:embed testdata/invalid-nested.yaml
var invalidNested []byte

func TestUnmarshalStrict(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  string
		data []byte
	}{
		{
			name: "valid",
			data: valid,
		},
		{
			name: "invalid",
			data: invalid,
			err:  "unknown keys",
		},
		{
			name: "invalid nested",
			data: invalidNested,
			err:  "this",
		},
		{
			name: "empty",
			data: []byte{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var a A

			err := xyaml.UnmarshalStrict(tt.data, &a)

			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)

				return
			}

			if len(tt.data) != 0 {
				require.NotEmpty(t, a)
			}

			require.NoError(t, err)
		})
	}
}
