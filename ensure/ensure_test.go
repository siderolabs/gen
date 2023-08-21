// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ensure_test

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/ensure"
)

func TestNoError(t *testing.T) {
	require.Panics(t, func() { ensure.NoError(errors.New("test error")) })
	require.NotPanics(t, func() { ensure.NoError(nil) })
}

func TestValue(t *testing.T) {
	require.Panics(t, func() { ensure.Value(net.ParseMAC("--02-00-5e-10-00-00-00-01")) })
	require.NotPanics(t, func() {
		require.Equal(t, "02:00:5e:10:00:00:00:01", ensure.Value(net.ParseMAC("02-00-5e-10-00-00-00-01")).String())
	})
}

func TestValues(t *testing.T) {
	require.Panics(t, func() { ensure.Values(net.ParseCIDR("192.-0.2.1--/24")) })
	require.NotPanics(t, func() {
		i, n := ensure.Values(net.ParseCIDR("192.0.2.1/24"))

		require.Equal(t, "192.0.2.1", i.String())
		require.Equal(t, "192.0.2.0/24", n.String())
	})
}
