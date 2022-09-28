// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xerrors_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xerrors"
)

//nolint:unused
type testTag struct{}

func TestTagged(t *testing.T) {
	stdErr := io.EOF
	taggedErr := xerrors.NewTagged[testTag](stdErr)

	require.True(t, xerrors.TypeIs[*xerrors.Tagged[testTag]](taggedErr))
	require.True(t, xerrors.TagIs[testTag](taggedErr))
	require.False(t, xerrors.TypeIs[*xerrors.Tagged[testTag]](stdErr))
	require.False(t, xerrors.TagIs[testTag](stdErr))
	require.ErrorIs(t, taggedErr, stdErr)

	taggedfErr := xerrors.NewTaggedf[testTag]("my custom error around: %w", stdErr)

	require.True(t, xerrors.TypeIs[*xerrors.Tagged[testTag]](taggedfErr))
	require.True(t, xerrors.TagIs[testTag](taggedfErr))
	require.ErrorIs(t, taggedfErr, stdErr)
	require.EqualError(t, taggedErr, stdErr.Error())

	taggedCustomErr := xerrors.NewTaggedf[testTag]("my custom error")

	require.True(t, xerrors.TypeIs[*xerrors.Tagged[testTag]](taggedCustomErr))
	require.True(t, xerrors.TagIs[testTag](taggedCustomErr))
	require.False(t, errors.Is(taggedCustomErr, stdErr))
}
