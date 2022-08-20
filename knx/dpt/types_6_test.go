package dpt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDPT_6010(t *testing.T) {
	knxValue := []byte{0, 42}
	dptValue := DPT_6010(42)

	var tmpDPT DPT_6010
	assert.NoError(t, tmpDPT.Unpack(knxValue))
	assert.Equal(t, dptValue, tmpDPT)

	assert.Equal(t, knxValue, dptValue.Pack())

	assert.Equal(t, "42 counter pulses", dptValue.String())
}
