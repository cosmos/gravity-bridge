package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalEthereumAddress(t *testing.T) {
	specs := map[string]struct {
		src EthereumAddress
		exp string
	}{
		"all good": {
			src: NewEthereumAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6"),
			exp: `"0xc783df8a850f42e7F7e57013759C285caa701eB6"`,
		},
		"empty address": {
			src: NewEthereumAddress(""),
			exp: `""`,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := spec.src.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, spec.exp, string(got))
		})
	}

}

func TestUnmarshalEthereumAddress(t *testing.T) {
	specs := map[string]struct {
		src      string
		expValue EthereumAddress
		expErr   bool
	}{
		"all good": {
			src:      `"0xc783df8a850f42e7F7e57013759C285caa701eB6"`,
			expValue: NewEthereumAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6"),
		},
		"empty address": {
			src:    `""`,
		},
		"non prefixed address": {
			src:    `"c783df8a850f42e7F7e57013759C285caa701eB6"`,
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			var got EthereumAddress
			err := got.UnmarshalJSON([]byte(spec.src))
			if spec.expErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, spec.expValue, got)
		})
	}
}
