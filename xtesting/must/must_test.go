// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package must_test

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xtesting/must"
)

func TestValue(t *testing.T) {
	value := must.Value(io.WriteString(io.Discard, "Hello, World!"))(t)
	require.Equal(t, 13, value)
}

func TestValues(t *testing.T) {
	i, n := must.Values(net.ParseCIDR("192.0.2.1/24"))(t)

	require.Equal(t, "192.0.2.1", i.String())
	require.Equal(t, "192.0.2.0/24", n.String())
}
