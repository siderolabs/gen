// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package check_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xtesting"
	"github.com/siderolabs/gen/xtesting/check"
)

func TestChecker(t *testing.T) {
	//nolint:govet
	tests := map[string]struct {
		arg      string
		expected int
		check    func(t xtesting.T, err error)
	}{
		"no error": {
			arg:      "1",
			expected: 1,
			check:    check.NoError(),
		},
		"error": {
			arg:      "a",
			expected: 0,
			check:    check.EqualError("strconv.Atoi: parsing \"a\": invalid syntax"),
		},
		"error contains": {
			arg:      "a",
			expected: 0,
			check:    check.ErrorContains("invalid syntax"),
		},
		"error regexp": {
			arg:      "a",
			expected: 0,
			check:    check.ErrorRegexp("strconv.*invalid syntax"),
		},
		"error as": {
			arg:      "a",
			expected: 0,
			check:    check.ErrorAs[*strconv.NumError](),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			atoi, err := strconv.Atoi(test.arg)
			test.check(t, err)

			require.Equal(t, test.expected, atoi)
		})
	}
}
