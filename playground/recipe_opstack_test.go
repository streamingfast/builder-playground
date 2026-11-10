package playground

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpRecipeValidate(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "defaults", args: nil},
		{name: "op-reth sequencer", args: []string{"--sequencer-el", "op-reth"}},
		{name: "karst after jovian", args: []string{"--sequencer-el", "op-reth", "--enable-latest-fork", "0", "--enable-karst", "10"}},
		{name: "unknown sequencer", args: []string{"--sequencer-el", "nethermind"}, wantErr: "invalid --sequencer-el"},
		{name: "karst on op-geth", args: []string{"--enable-latest-fork", "0", "--enable-karst", "10"}, wantErr: "op-geth does not support Karst"},
		{name: "karst at genesis", args: []string{"--sequencer-el", "op-reth", "--enable-latest-fork", "0", "--enable-karst", "0"}, wantErr: "after genesis"},
		{name: "karst without jovian", args: []string{"--sequencer-el", "op-reth", "--enable-karst", "10"}, wantErr: "requires Jovian"},
		{name: "karst before jovian", args: []string{"--sequencer-el", "op-reth", "--enable-latest-fork", "10", "--enable-karst", "10"}, wantErr: "requires Jovian"},
		{name: "op-reth twice", args: []string{"--sequencer-el", "op-reth", "--external-builder", "op-reth"}, wantErr: "cannot be used with"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recipe := &OpRecipe{}
			require.NoError(t, recipe.Flags().Parse(tc.args))

			err := recipe.Validate()
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestOpRecipeSequencerEL(t *testing.T) {
	for _, sequencerEL := range []string{"op-geth", "op-reth"} {
		t.Run(sequencerEL, func(t *testing.T) {
			recipe := &OpRecipe{}
			require.NoError(t, recipe.Flags().Parse([]string{"--sequencer-el", sequencerEL}))

			o := newTestOutput(t)
			manifest := NewManifest("test", recipe.Apply(&ExContext{Output: o, Contender: &ContenderContext{}}))

			_, ok := manifest.GetService(sequencerEL)
			require.True(t, ok, "missing %s service", sequencerEL)

			opNode, ok := manifest.GetService("op-node")
			require.True(t, ok)
			require.Contains(t, opNode.Args, Connect(sequencerEL, "authrpc"))
		})
	}
}
