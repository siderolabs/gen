// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package channel

import "context"

// RecvWithContext tries to receive a value from a channel. Operation is aborted if the context is canceled.
//
// Function returns [StateRecv] if the value was received, [StateCancelled] on context cancelation and
// [StateClosed] if the channel was closed.
func RecvWithContext[T any](ctx context.Context, ch <-chan T) (T, RecvState) {
	select {
	case <-ctx.Done():
		var zero T

		return zero, StateCancelled
	case val, ok := <-ch:
		if !ok {
			return val, StateClosed
		}

		return val, StateRecv
	}
}

// RecvState is the state of a channel after receiving a value.
type RecvState int

const (
	// StateRecv means that a value was received from the channel.
	StateRecv RecvState = iota
	// StateEmpty means that the channel was empty.
	StateEmpty
	// StateClosed means that the channel was closed.
	StateClosed
	// StateCancelled means that the context was canceled.
	StateCancelled
)

// TryRecv tries to receive a value from a channel.
//
// Function returns the value and the state of the channel.
func TryRecv[T any](ch <-chan T) (T, RecvState) {
	var zero T

	select {
	case val, ok := <-ch:
		if !ok {
			return zero, StateClosed
		}

		return val, StateRecv
	default:
		return zero, StateEmpty
	}
}
