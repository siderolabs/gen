// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ordered_test

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/pair/ordered"
)

func TestTriple(t *testing.T) {
	t.Parallel()

	expectedSlice := []ordered.Triple[int, string, float64]{
		ordered.MakeTriple(math.MinInt64, "Alpha", 69.0),
		ordered.MakeTriple(-200, "Alpha", 69.0),
		ordered.MakeTriple(-200, "Beta", -69.0),
		ordered.MakeTriple(-200, "Beta", 69.0),
		ordered.MakeTriple(1, "", 69.0),
		ordered.MakeTriple(1, "Alpha", 67.0),
		ordered.MakeTriple(1, "Alpha", 68.0),
		ordered.MakeTriple(10, "Alpha", 68.0),
		ordered.MakeTriple(10, "Beta", 68.0),
		ordered.MakeTriple(math.MaxInt64, "", 69.0),
	}

	seed1 := time.Now().UnixNano()
	seed2 := time.Now().UnixNano()
	rnd := rand.New(rand.NewPCG(uint64(seed1), uint64(seed2)))

	for i := range 1000 {
		a := append([]ordered.Triple[int, string, float64](nil), expectedSlice...)
		rnd.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

		slices.SortFunc(a, func(i, j ordered.Triple[int, string, float64]) int { return i.Compare(j) })
		require.Equal(t, expectedSlice, a, "failed with seed1 %d seed2 %d iteration %d", seed1, seed2, i)
	}
}
