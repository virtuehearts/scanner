package brainforce

import (
	"testing"
	"time"

	"github.com/shlima/fortune/internal/pkg/datum"
	"github.com/shlima/fortune/internal/pkg/key"
	"github.com/shlima/fortune/internal/pkg/pass"
	"github.com/stretchr/testify/require"
)

func TestForce_Generate(t *testing.T) {
	t.Parallel()

	t.Run("when nothing found", func(t *testing.T) {
		t.Parallel()
		setup := MustSetup(t)
		defer setup.ctrl.Finish()

		setup.pass.EXPECT().
			Next().
			Return("", pass.ErrEnd)

		out := make([]key.KeyChain, 0)
		err := setup.Generate(func(chain key.KeyChain) {
			out = append(out, chain)
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(out))
	})

	t.Run("when found", func(t *testing.T) {
		t.Parallel()
		setup := MustSetup(t)
		defer setup.ctrl.Finish()

		chain := key.KeyChain{
			Private:      "Private",
			Compressed:   "Compressed",
			Uncompressed: "Uncompressed",
		}

		switch {
		case time.Now().Second()%2 == 0:
			setup.SetIndex(datum.NewFilteredIndex(datum.Index{chain.Compressed: true}, nil))
		default:
			setup.SetIndex(datum.NewFilteredIndex(datum.Index{chain.Uncompressed: true}, nil))
		}

		first := setup.pass.EXPECT().
			Next().
			Return("foo", nil)

		setup.pass.EXPECT().
			Next().
			Return("", pass.ErrEnd).
			After(first)

		setup.key.EXPECT().
			BrainSHA256([]byte("foo")).
			Return(chain, nil)

		out := make([]key.KeyChain, 0)
		err := setup.Generate(func(chain key.KeyChain) {
			out = append(out, chain)
		})

		require.NoError(t, err)
		require.Equal(t, 1, len(out))
		require.Contains(t, out, chain)
	})
}
