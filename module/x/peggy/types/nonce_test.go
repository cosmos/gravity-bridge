package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonceGreaterThan(t *testing.T) {
	specs := map[string]struct {
		src   UInt64Nonce
		other nonce
		exp   bool
	}{
		"equal":   {NewUInt64Nonce(1), NewUInt64Nonce(1), false},
		"greater": {NewUInt64Nonce(2), NewUInt64Nonce(1), true},
		"less":    {NewUInt64Nonce(1), NewUInt64Nonce(2), false},
		//"src nil":     {nil, NewUInt64Nonce(2), false},
		"other nil":         {NewUInt64Nonce(1), nil, true},
		"other nil pointer": {NewUInt64Nonce(1), (*UInt64Nonce)(nil), true},
		//"both nil":    {nil, nil, false},
		"src empty":   {NewUInt64Nonce(0), NewUInt64Nonce(2), false},
		"other empty": {NewUInt64Nonce(1), NewUInt64Nonce(0), true},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got := spec.src.GreaterThan(spec.other)
			assert.Equal(t, spec.exp, got)
		})
	}
}
