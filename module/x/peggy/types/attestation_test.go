package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonceGreaterThan(t *testing.T) {
	specs := map[string]struct {
		src, other Nonce
		exp        bool
	}{
		"equal":       {NonceFromUint64(1), NonceFromUint64(1), false},
		"greater":     {NonceFromUint64(2), NonceFromUint64(1), true},
		"less":        {NonceFromUint64(1), NonceFromUint64(2), false},
		"src nil":     {nil, NonceFromUint64(2), false},
		"other nil":   {NonceFromUint64(1), nil, true},
		"src empty":   {Nonce{}, NonceFromUint64(2), false},
		"other empty": {NonceFromUint64(1), Nonce{}, true},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got := spec.src.GreaterThan(spec.other)
			assert.Equal(t, spec.exp, got)
		})
	}

}
