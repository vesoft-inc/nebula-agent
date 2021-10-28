package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	assert := assert.New(t)
	addrStr := "127.0.0.1:80"

	addr, err := ParseAddr(addrStr)
	assert.Nil(err)
	assert.Equal(addr.Host, "127.0.0.1")
	assert.Equal(addr.Port, int32(80))

	encStr := StringifyAddr(addr)
	assert.Equal(addrStr, encStr)
}
