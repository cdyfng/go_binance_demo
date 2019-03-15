package util

import (
	"github.com/cdyfng/go_binance_demo/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtil(t *testing.T) {

	f := 0.0000000111
	r := 0.00000001
	time8_r := int64(1)
	assert.Equal(t, util.Round8(f), r)
	assert.Equal(t, util.Times8(f), time8_r)
}
