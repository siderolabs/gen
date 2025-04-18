// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package channel_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/siderolabs/gen/channel"
)

func TestSendWithContext(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 1)
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	assert.True(t, channel.SendWithContext(ctx, ch, 42))
	assert.Equal(t, 42, <-ch)

	assert.True(t, channel.SendWithContext(ctx, ch, 69))

	cancel()
	assert.False(t, channel.SendWithContext(ctx, ch, 33))
}
