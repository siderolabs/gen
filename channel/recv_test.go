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

func TestRecvWithContext(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 1)
	ch <- 42

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	val, state := channel.RecvWithContext(ctx, ch)
	assert.Equal(t, 42, val)
	assert.Equal(t, channel.StateRecv, state)

	cancel()

	val, state = channel.RecvWithContext(ctx, ch)
	assert.Zero(t, val)
	assert.Equal(t, channel.StateCancelled, state)
}

func TestRecvWithContextCloseCh(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 1)
	ch <- 42

	ctx := context.Background()

	val, state := channel.RecvWithContext(ctx, ch)
	assert.Equal(t, 42, val)
	assert.Equal(t, channel.StateRecv, state)

	close(ch)
	val, state = channel.RecvWithContext(ctx, ch)
	assert.Zero(t, val)
	assert.Equal(t, channel.StateClosed, state)
}

func TestTryRecv(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 1)
	ch <- 42

	val, state := channel.TryRecv(ch)
	assert.Equal(t, 42, val)
	assert.Equal(t, channel.StateRecv, state)

	val, state = channel.TryRecv(ch)
	assert.Zero(t, val)
	assert.Equal(t, channel.StateEmpty, state)

	close(ch)
	val, state = channel.TryRecv(ch)
	assert.Zero(t, val)
	assert.Equal(t, channel.StateClosed, state)
}
