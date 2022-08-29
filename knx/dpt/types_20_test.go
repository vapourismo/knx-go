package dpt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDPT_20102(t *testing.T) {
	knxValue := []byte{0, 4}
	dptValue := DPT_20102(4)

	var tmpDPT DPT_20102
	assert.NoError(t, tmpDPT.Unpack(knxValue))
	assert.Equal(t, dptValue, tmpDPT)

	assert.Equal(t, knxValue, dptValue.Pack())

	assert.Equal(t, "Building Protection", dptValue.String())
}
